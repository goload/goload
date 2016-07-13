import React from "react";
import {Glyphicon, OverlayTrigger, Tooltip} from "react-bootstrap";

export class ExtractIndicator extends React.Component {
    constructor(props) {
        super(props);
    }

    render() {
        return (<OverlayTrigger placement="top" overlay={<Tooltip id="extractingTooltip">Extracting</Tooltip>}>
            <Glyphicon
                glyph="cog" className="gly-spin"/>
        </OverlayTrigger>)
    }
}