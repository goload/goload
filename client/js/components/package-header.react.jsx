import React from "react";
import {ProgressBar, Row, Col, Glyphicon, OverlayTrigger, Tooltip} from "react-bootstrap";

export class PackageHeader extends React.Component {
    constructor(props) {
        super(props)
    }

    render() {
        return (<div>
            <Row >
                <Col sm={6} role="button"
                     onClick={this.props.toggleCollapse}>
                    <OverlayTrigger placement="top" overlay={<Tooltip id="expandTooltip">Expand</Tooltip>}>
                        <Glyphicon
                            glyph="folder-open"/>
                    </OverlayTrigger>
                    {' '}{' '}{this.props.package.name}{' '}{'(' + formatBytes(this.props.package.size, 1) + ')'}
                </Col>

                <Col className="text-right" smOffset={5} sm={1}>
                    <OverlayTrigger placement="top" overlay={<Tooltip id="retryTooltip">Retry</Tooltip>}>
                        <Glyphicon role="button"
                                   onClick={this.props.retryPackage}
                                   glyph="refresh"/>
                    </OverlayTrigger>
                    {' '}
                    <OverlayTrigger placement="top" overlay={<Tooltip id="removeTooltip">Remove</Tooltip>}>
                        <Glyphicon role="button"
                                   onClick={this.props.removePackage}
                                   glyph="trash"/>
                    </OverlayTrigger>
                </Col>
            </Row>
            <Row>
                <Col sm={12}>
                    <ProgressBar bsStyle="success" active={this.props.package.progress < 100}
                                 now={this.props.package.progress}
                                 label={Math.round(this.props.package.progress)+"%"}/>
                </Col>
            </Row>
        </div>)
    }
}
function formatBytes(bytes, decimals) {
    if (bytes == 0) return '0 Byte';
    var k = 1000; // or 1024 for binary
    var dm = decimals + 1 || 3;
    var sizes = ['Bytes', 'KB', 'MB', 'GB', 'TB', 'PB', 'EB', 'ZB', 'YB'];
    var i = Math.floor(Math.log(bytes) / Math.log(k));
    return parseFloat((bytes / Math.pow(k, i)).toFixed(dm)) + ' ' + sizes[i];
}