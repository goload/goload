import React from 'react'
import {Glyphicon, Navbar, Nav, NavItem} from 'react-bootstrap';
import {LinkContainer} from 'react-router-bootstrap';
import {Link} from 'react-router'


export class NavBar extends React.Component {
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
                            <LinkContainer to={'settings'}><NavItem > <Glyphicon glyph="wrench"/>{' '}Settings</NavItem></LinkContainer>
                        </Nav>
                    </Navbar.Collapse>
                </Navbar>
            </div>
        )
    }
}
