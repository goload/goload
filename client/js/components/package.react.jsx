import React from 'react'
import {ProgressBar, Collapse, Row, Col, Glyphicon, OverlayTrigger, Tooltip} from 'react-bootstrap'
import moment from 'moment'
require("moment-duration-format");
export class Package extends React.Component {
    constructor(props) {
        super(props)
        this.state = {
            collapsed: true
        }
    }

    render() {
        return (
            <Row >
                <Col sm={10}>
                    <Row >
                        <Col sm={6} role="button"
                             onClick={() => this.setState({collapsed: !this.state.collapsed}) }>
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
                    <Collapse in={!this.state.collapsed}>
                        <div>
                            {this.props.package.files.map(file=> {
                                    var barStyle = null;
                                    if (file.failed) {
                                        barStyle = "danger"
                                    }
                                    var glyph = "save-file";
                                    if (file.finished) {
                                        if (!file.failed) {
                                            glyph = "saved";
                                        } else {
                                            glyph = "remove-circle"
                                        }
                                    }
                                    return (<Row key={file.url}>

                                        {/*<Col sm={1}><Glyphicon glyph={glyph} style={{'fontSize':'2.2em'}}/></Col>*/}
                                        <Col smOffset={1} sm={11}>
                                            <Row>
                                                <Col sm={6}><Glyphicon
                                                    glyph="compressed"/>{' '}{file.filename != '' ? file.filename : file.url}{' '}
                                                    ({formatBytes(file.size)})</Col>
                                                <Col className="text-right" sm={6}>
                                                    <Glyphicon
                                                        glyph="time"/> {moment.duration(file.ete/(1000*1000)).format()}
                                                    {' '}
                                                    <Glyphicon
                                                        glyph="arrow-down"/> {file.download_speed}</Col>
                                            </Row>
                                            <ProgressBar className="progress-file" bsStyle={barStyle} active={file.progress < 100}
                                                              now={file.progress}
                                                              label={Math.round(file.progress)+"%"}/>

                                        </Col>
                                    </Row>)
                                }
                            )}
                        </div>
                    </Collapse>
                </Col>
            </Row>
        )
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