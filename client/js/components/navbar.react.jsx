import React from 'react'
import {Glyphicon, Navbar, Nav, NavItem} from 'react-bootstrap';
import {LinkContainer} from 'react-router-bootstrap';
import {Link} from 'react-router'
import $ from 'jquery'

export class NavBar extends React.Component {
    constructor(props){
        super(props);
        this.state={
            version:''
        }
    }

    componentDidMount() {
        $.get('/api/version',(json)=>
            this.setState({version:json.version})
        )
    }

    render() {
        return (
            <div>
                <Navbar >
                    <Navbar.Header>
                        <Navbar.Brand>
                            <Link to={'/'}><Glyphicon glyph="download"/>{' '}Uploaded Downloader</Link>
                        </Navbar.Brand>
                        <Navbar.Toggle/>
                    </Navbar.Header>
                    <Navbar.Collapse>
                        <Nav>
                            <LinkContainer to={'/'}><NavItem > <Glyphicon
                                glyph="home"/>{' '}Home</NavItem></LinkContainer>
                        </Nav>
                        <Nav pullRight>
                            <NavItem target="_blank" href={'https://github.com/goload/goload/releases/tag/v'+ this.state.version}>{this.state.version}</NavItem>
                        </Nav>
                        <Nav pullRight>
                            <LinkContainer to={'settings'}><NavItem > <Glyphicon glyph="wrench"/>{' '}Settings</NavItem></LinkContainer>
                        </Nav>
                    </Navbar.Collapse>
                </Navbar>
            </div>
        )
    }
}
