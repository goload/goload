import React from "react";
import {
    Col,
    Tooltip,
    OverlayTrigger,
    FormGroup,
    FormControl,
    Glyphicon,
    ControlLabel,
    Form,
    Button
} from "react-bootstrap";
import {Package} from "./package.react.jsx";

import $ from "jquery";
import _ from "lodash";
import alertify from "alertify.js";
alertify.logPosition("top right");

const _url = '/api/packages';

export class PackageView extends React.Component {
    constructor(props) {
        super(props);
        this.state = {
            packages: [],
            packageName: '',
            links: '',
            password: ''
        };
        this.handleLinks = this.handleLinks.bind(this);
        this.handlePassword = this.handlePassword.bind(this);
        this.handlePackageName = this.handlePackageName.bind(this);
        this.submitPackage = this.submitPackage.bind(this);
        this.loadPackages = this.loadPackages.bind(this);
    }

    componentDidMount() {
        this.startPolling()
    }

    componentWillUnmount() {
        if (this._timer) {
            clearInterval(this._timer);
            this._timer = null;
        }
    }

    startPolling() {
        this.loadPackages();
        this._timer = setInterval(() => {
            this.loadPackages();
        }, 3000);
    }

    loadPackages() {
        $.get(_url, result => {
            this.setState({
                packages: _.sortBy(result, (item)=>item.date_added).reverse()
            });
        });
    }

    handlePassword(e) {
        this.setState({password: e.target.value});
    }

    handlePackageName(e) {
        this.setState({packageName: e.target.value});
    }

    handleLinks(e) {
        this.setState({links: e.target.value});
    }

    submitPackage() {
        var splittetLinks = this.state.links.trim().split(/\s/g);
        var data = {
            name: this.state.packageName,
            files: [],
            password: this.state.password
        };
        _.forEach(splittetLinks, link => {
            if (link !== '') {
                data.files.push({'url': link})
            }
        });
        $.post(_url, JSON.stringify(data)).done(()=> {

            alertify.delay(2000).success('Package added');
            this.setState({
                packageName: '',
                links: '',
                password: ''
            });
            this.loadPackages()
        }).fail(()=> {
            alertify.delay(2000).error('No package name provided');
        })

    }

    retryPackage(pack) {
        $.get(_url+'/' + pack.id + '/retry', ()=> {
                this.loadPackages();
                alertify.delay(2000).success('Retrying package ' + pack.name);
            }
        );
    }

    removePackage(pack) {
        $.ajax({
            url: _url+'/' + pack.id,
            type: 'DELETE',
            success: ()=> {
                this.loadPackages();
                alertify.delay(2000).success('Package ' + pack.name + ' removed');
            }
        });
    }

    render() {
        const tooltip = (
            <Tooltip id="tooltip"><strong>Only uploaded.to etc. links</strong> Please decrypt dlc files yourself for
                now.</Tooltip>
        );
        return (
            <div>
                <Form horizontal>
                    <FormGroup >
                        <Col componentClass={ControlLabel} sm={2}>
                            Package Name
                        </Col>
                        <Col sm={4}>
                            <FormControl col
                                         type="text"
                                         value={this.state.packageName}
                                         placeholder="Package Name"
                                         onChange={this.handlePackageName}/>
                        </Col>
                    </FormGroup>


                    <FormGroup >
                        <Col componentClass={ControlLabel} sm={2}>
                            Extract Password
                        </Col>
                        <Col sm={4}>
                            <FormControl col
                                         type="text"
                                         value={this.state.password}
                                         placeholder="Password"
                                         onChange={this.handlePassword}/>
                        </Col>
                    </FormGroup>
                    <FormGroup >
                        <Col componentClass={ControlLabel} sm={2}>
                            Links{' '}
                            <OverlayTrigger placement="right" overlay={tooltip}>
                                <a target="_blank" href="http://dcrypt.it"> <Glyphicon glyph="info-sign"/></a>
                            </OverlayTrigger>
                        </Col>
                        <Col sm={4}>
                            <FormControl col
                                         type="text"
                                         value={this.state.links}
                                         placeholder="Links"
                                         onChange={this.handleLinks}/>
                        </Col>
                    </FormGroup>

                    <FormGroup>
                        <Col smOffset={2} sm={10}>
                            <Button onClick={this.submitPackage}><Glyphicon glyph="plus"/> Add package</Button>
                        </Col>
                    </FormGroup>
                </Form>
                {this.state.packages.map(pack =>
                    <Package key={pack.id} package={pack} removePackage={this.removePackage.bind(this,pack)} retryPackage={this.retryPackage.bind(this,pack)}/>
                )}
            </div>
        )
    }
}


