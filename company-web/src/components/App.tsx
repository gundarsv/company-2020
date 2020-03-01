import * as React from "react";
import {BrowserRouter as Router, Route} from "react-router-dom";
import Company from "./company/company";
import Navigation from "./navigation/navigation";
import Home from "./home/home";
import {SnackbarProvider} from 'notistack';
import Owner from "./owner/owner";

const App: React.FC = () => {
	return (
		<Router>
			<SnackbarProvider maxSnack={3}>
				<Navigation/>
				<Route exact={true} path="/" component={Home}/>
				<Route exact={true} path="/company" component={Company}/>
				<Route exact={true} path="/owner" component={Owner}/>
			</SnackbarProvider>
		</Router>
	);
};

export default App;
