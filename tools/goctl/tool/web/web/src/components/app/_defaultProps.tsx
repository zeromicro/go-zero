import React, {CSSProperties} from "react";
import {ConverterIcon} from "../../util/icon";
import type {TFunction} from "i18next";

const iconStyle: CSSProperties = {
    color: "#575757",
    fontSize: 24,
    textAlign: "center",
    marginLeft: "-4px"
}

const subMenuIconStyle: CSSProperties = {
    color: "#575757",
    fontSize: 18,
    textAlign: "center",
    marginLeft: "-4px"
}

export const menuItems = (t: TFunction) => {
    return [
        {
            key: 'welcome',
            icon: <ConverterIcon style={iconStyle} type={"icon-welcome"}/>,
            label: t("menuWelcome")
        },
        {
            key: 'api',
            icon: <ConverterIcon style={iconStyle} type={"icon-code"}/>,
            label: t("api"),
            children: [
                {
                    key: 'builder',
                    icon: <ConverterIcon style={subMenuIconStyle} type={"icon-builder"}/>,
                    label: t("builder"),
                },
            ]
        },
    ]
}