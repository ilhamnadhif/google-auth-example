import "./App.css";
import { GoogleLogin } from 'react-google-login';

function App() {

  const responseGoogle = (response) =>{
    console.log(response);
  }

  return (
    <div className="App">
      <GoogleLogin
        clientId="351570360340-c2fh8e6t265el6d73a0n25m95s9uq49j.apps.googleusercontent.com"
        buttonText="Login"
        onSuccess={responseGoogle}
        onFailure={responseGoogle}
        cookiePolicy={"single_host_origin"}
      />
    </div>
  );
}

export default App;
