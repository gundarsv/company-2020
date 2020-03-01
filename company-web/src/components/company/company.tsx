import * as React from "react";
import { withRouter, useHistory } from "react-router-dom";


const Company: React.FC = () => {
  const history = useHistory();
  
  return (
    <div>
      <p>Company</p>
    </div>
  );
};

export default withRouter(Company);
