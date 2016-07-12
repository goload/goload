import React from "react";
import {Collapse, Row, Col} from "react-bootstrap";
import {File} from "./file.react.jsx";
import {PackageHeader} from "./package-header.react.jsx";
require("moment-duration-format");
export class Package extends React.Component {
    constructor(props) {
        super(props);
        this.state = {
            collapsed: true
        }
    }

    render() {
        return (
            <Row >
                <Col sm={12}>
                    <PackageHeader removePackage={this.props.removePackage}
                                   retryPackage={this.props.retryPackage}
                                   toggleCollapse={() => this.setState({collapsed: !this.state.collapsed}) }
                                   package={this.props.package}/>

                    <Collapse in={!this.state.collapsed}>
                        <div>
                            {this.props.package.files.map((file, index)=>
                                <File key={index} file={file}/>
                            )}
                        </div>
                    </Collapse>
                </Col>
            </Row>
        )
    }
}

