import React from "react";
import {ProgressBar, Row, Col, Glyphicon, OverlayTrigger, Tooltip} from "react-bootstrap";
import {ExtractIndicator} from './extract-indicator.react.jsx'
export class PackageHeader extends React.Component {
    constructor(props) {
        super(props)
    }

    render() {
        var progress = this.props.package.progress;
        var indicator;
        if (this.props.package.extracting) {
            progress = this.props.package.unrar_progress;
            indicator = (<ExtractIndicator/>);
        }
        return (<div>
            <Row >
                <Col sm={6} role="button"
                     onClick={this.props.toggleCollapse}>
                    {indicator}{' '}
                    <Glyphicon
                        glyph="folder-open"/>
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
                    <ProgressBar bsStyle="success" active={progress < 100}
                                 now={progress}
                                 label={Math.round(progress)+"%"}/>
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