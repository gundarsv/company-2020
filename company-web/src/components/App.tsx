import * as React from "react";
import { BrowserRouter as Router, Route } from "react-router-dom";
import Company from "./company/company";
import Navigation from "./navigation/navigation";
import Home from "./home/home";

const App: React.FC = () => {
  return (
      <Router>
        <Navigation />
        <Route exact={true} path="/" component={Home} />
        <Route exact={true} path="/company" component={Company} />
        <Route exact={true} path="/owner" component={Company} />
      </Router>
  );
};

export default App;
