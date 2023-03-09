import SwaggerUI from "swagger-ui-react";
import "swagger-ui-react/swagger-ui.css"
import { proto } from "./proto";

function App() {

  return (
    <SwaggerUI spec={proto.replace("__SERVER_ADDR__", window.location.origin)} />
  );
}

export default App
