import "./App.css";
import { GoogleLogin } from "react-google-login";
import axios from "axios";
import googleOneTap from "google-one-tap";
import { useState } from "react";

const googleClientId =
  "351570360340-c2fh8e6t265el6d73a0n25m95s9uq49j.apps.googleusercontent.com";

function App() {
  const [user, setUser] = useState({});

  const successHandler = async (response) => {
    try {
      const res = await axios.post(
        "http://localhost:5000/google",
        JSON.stringify({ token: response.tokenId }),
        {
          headers: {
            "Content-Type": "application/json",
          },
        }
      );
      const resUser = await axios.get("http://localhost:5000/me", {
        headers: {
          Authorization: `Bearer ${res.data.token}`,
        },
      });

      setUser(resUser.data);
    } catch (error) {
      setUser({});
      console.log(error);
    }
  };
  const successOneTapHandler = async (response) => {
    try {
      const res = await axios.post(
        "http://localhost:5000/google",
        JSON.stringify({ token: response.credential }),
        {
          headers: {
            "Content-Type": "application/json",
          },
        }
      );
      const resUser = await axios.get("http://localhost:5000/me", {
        headers: {
          Authorization: `Bearer ${res.data.token}`,
        },
      });

      setUser(resUser.data);
    } catch (error) {
      setUser({});
      console.log(error);
    }
  };
  const failureHandler = (response) => {
    console.log(response);
  };
  googleOneTap(
    {
      client_id: googleClientId,
      auto_select: false,
      cancel_on_tap_outside: false,
      context: "signin",
    },
    successOneTapHandler
  );

  return (
    <div className="App">
      <GoogleLogin
        clientId={googleClientId}
        buttonText="Login"
        onSuccess={successHandler}
        onFailure={failureHandler}
        cookiePolicy={"single_host_origin"}
      />
      <button onClick={() => setUser({})}>Logout</button>
      {user.email && <h1>{user.email}</h1>}
    </div>
  );
}

export default App;
