import "./App.css";
import { GoogleLogin } from "react-google-login";
import axios from "axios";
import { useState } from "react";

function App() {
  const [user, setUser] = useState({})
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

      setUser(resUser.data)
      
    } catch (error) {
      setUser({})
      console.log(error);
    }
  };
  const failureHandler = (response) => {
    alert(response);
  };


  return (
    <div className="App">
      <GoogleLogin
        clientId="351570360340-c2fh8e6t265el6d73a0n25m95s9uq49j.apps.googleusercontent.com"
        buttonText="Login"
        onSuccess={successHandler}
        onFailure={failureHandler}
        cookiePolicy={"single_host_origin"}
      />
      <button onClick={()=>setUser({})}>Logout</button>
      {user.email && <h1>{user.email}</h1>}
    </div>
  );
}

export default App;
