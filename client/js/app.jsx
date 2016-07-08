import React from "react";
import {render} from "react-dom";
import {Router, Route, browserHistory} from "react-router";
import {Configuration} from "./components/config.react.jsx";
import {NavBar} from './components/navbar.react.jsx'
import {PackageView} from './components/package-view.react.jsx'

class App extends React.Component{
    render() {
        return (
            <div>
                <NavBar />
                <div className="container">
                    {this.props.children}
                </div>
            </div>
        )
    }
}


render((
    <Router history={browserHistory}>
        <Route path="" component={App}>
            <Route path="/" component={PackageView}/>
            <Route path="settings" component={Configuration}/>
        </Route>
    </Router>
), document.getElementById('root'));