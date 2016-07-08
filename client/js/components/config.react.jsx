import React from "react";
import {Form, Col, FormGroup, FormControl, Glyphicon, ControlLabel, Button} from "react-bootstrap";
import $ from 'jquery'
import alertify from "alertify.js";
alertify.logPosition("top right");
const _url = '/config';
export class Configuration extends React.Component {

    constructor(props) {
        super(props);
        this.state = {
            account: {
                username: '',
                password: ''
            },
            dirs: {
                downloadDir: '',
                extractDir: ''
            }

        };
        this.handleUsername = this.handleUsername.bind(this);
        this.handlePassword = this.handlePassword.bind(this);
        this.handleDownloadDir = this.handleDownloadDir.bind(this);
        this.handleExtractDir = this.handleExtractDir.bind(this);
        this.saveDirs = this.saveDirs.bind(this);
        this.saveAccount = this.saveAccount.bind(this);
        this.loadConfig = this.loadConfig.bind(this);
    }

    componentDidMount() {
        this.loadConfig();
    }

    loadConfig() {
        $.get(_url, result => {
            this.setState(result);
        })
    }

    handleUsername(e) {
        var newState = this.state;
        newState.account.username = e.target.value;
        this.setState(newState)
    }

    handlePassword(e) {
        var newState = this.state;
        newState.account.password = e.target.value;
        this.setState(newState)
    }

    handleDownloadDir(e) {
        var newState = this.state;
        newState.dirs.downloadDir = e.target.value;
        this.setState(newState)
    }

    handleExtractDir(e) {
        var newState = this.state;
        newState.dirs.extractDir = e.target.value;
        this.setState(newState)
    }

    saveAccount() {
        
        $.ajax({
            url:_url+"/account",
            type:'PUT',
            data:JSON.stringify(this.state.account),
            contentType: 'application/json; charset=utf-8',
            success: () => {
                alertify.delay(2000).success("Account saved");
            }
        })
    }

    saveDirs() {

        $.ajax({
            url:_url+"/dirs",
            type:'PUT',
            data:JSON.stringify(this.state.dirs),
            contentType: 'application/json; charset=utf-8',
            success: () => {
                alertify.delay(2000).success('Directories saved');
            }
        })
    }

    render() {
        return (<Form horizontal>
            <Col sm={2}/><h4>Account Information</h4>
            <FormGroup >
                <Col componentClass={ControlLabel} sm={2}>
                    Username
                </Col>
                <Col sm={4}>
                    <FormControl col
                                 type="text"
                                 value={this.state.account.username}
                                 placeholder="Username"
                                 onChange={this.handleUsername}/>
                </Col>
            </FormGroup>


            <FormGroup >
                <Col componentClass={ControlLabel} sm={2}>
                    Password
                </Col>
                <Col sm={4}>
                    <FormControl col
                                 type="password"
                                 value={this.state.account.password}
                                 placeholder="Password"
                                 onChange={this.handlePassword}/>
                </Col>
            </FormGroup>
            <FormGroup>
                <Col smOffset={2} sm={10}>
                    <Button onClick={this.saveAccount}><Glyphicon glyph="save"/>{' '}Save Account Information</Button>
                </Col>
            </FormGroup>
            <Col sm={2}/><h4 componentClass={ControlLabel}>Directories</h4>
            <FormGroup >
                <Col componentClass={ControlLabel} sm={2}>
                    Download Directory
                </Col>
                <Col sm={4}>
                    <FormControl col
                                 type="text"
                                 value={this.state.dirs.downloadDir}
                                 placeholder="Download Directory"
                                 onChange={this.handleDownloadDir}/>
                </Col>
            </FormGroup>

            <FormGroup >
                <Col componentClass={ControlLabel} sm={2}>
                    Extract Directory
                </Col>
                <Col sm={4}>
                    <FormControl col
                                 type="text"
                                 value={this.state.dirs.extractDir}
                                 placeholder="Extract Directory"
                                 onChange={this.handleExtractDir}/>
                </Col>
            </FormGroup>

            <FormGroup>
                <Col smOffset={2} sm={10}>
                    <Button onClick={this.saveDirs}><Glyphicon glyph="save"/>{' '}Save Directories</Button>
                </Col>
            </FormGroup>
        </Form>)
    }
}