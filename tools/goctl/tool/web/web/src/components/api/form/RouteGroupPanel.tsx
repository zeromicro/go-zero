import React from "react";
import { FormListFieldData } from "antd/es/form/FormList";
import RouteGroupOptionPanel from "./RouteGroupOptionPanel";
import RoutePanel from "./RoutePanel";
import type { FormInstance } from "antd/es/form/hooks/useForm";

interface RouteGroupPanelProps {
  routeGroupField: FormListFieldData;
  form: FormInstance;
}

const RouteGroupPanel: React.FC<
  RouteGroupPanelProps & React.RefAttributes<HTMLDivElement>
> = (props) => {
  const routeGroupField = props.routeGroupField;
  return (
    <div>
      <RouteGroupOptionPanel routeGroupField={routeGroupField} />
      <RoutePanel routeGroupField={routeGroupField} form={props.form} />
    </div>
  );
};

export default RouteGroupPanel;
