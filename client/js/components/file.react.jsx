/**
 * Created by andre on 12.07.16.
 */
import React from "react";
import {ProgressBar, Row, Col, Glyphicon, OverlayTrigger, Tooltip} from "react-bootstrap";
import moment from "moment";
export class File extends React.Component {
    constructor(props) {
        super(props)
    }

    render() {
        var file = this.props.file;
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
        var cog;
        if(file.extracting) {
            cog = (<OverlayTrigger placement="top" overlay={<Tooltip id="extractingTooltip">Extracting</Tooltip>}>
                <Glyphicon
                    glyph="cog" className="gly-spin"/>
            </OverlayTrigger>);
        }
        return (<Row>
            {/*<Col sm={1}><Glyphicon glyph={glyph} style={{'fontSize':'2.2em'}}/></Col>*/}
            <Col smOffset={0} sm={12}>
                <Row>
                    <Col sm={6}>
                        {cog}{' '}
                        <Glyphicon
                        glyph="compressed"/>
                        {' '}{file.filename != '' ? file.filename : file.url}{' '}
                        ({formatBytes(file.size)}){' '}
                    </Col>
                    <Col className="text-right" sm={6}>
                        <Glyphicon
                            glyph="time"/> {moment.duration(file.ete / (1000 * 1000)).format()}
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
}

function formatBytes(bytes, decimals) {
    if (bytes == 0) return '0 Byte';
    var k = 1000; // or 1024 for binary
    var dm = decimals + 1 || 3;
    var sizes = ['Bytes', 'KB', 'MB', 'GB', 'TB', 'PB', 'EB', 'ZB', 'YB'];
    var i = Math.floor(Math.log(bytes) / Math.log(k));
    return parseFloat((bytes / Math.pow(k, i)).toFixed(dm)) + ' ' + sizes[i];
}