import * as React from "react";
import { makeStyles } from "@material-ui/core/styles";
import Paper from "@material-ui/core/Paper";
import Tabs from "@material-ui/core/Tabs";
import Tab from "@material-ui/core/Tab";
import { withRouter, useHistory } from "react-router-dom";

const useStyles = makeStyles({
  root: {
    flexGrow: 1
  }
});

const Navigation: React.FC = () => {
    const classes = useStyles();
    const history = useHistory();

    const handleChange = (event: React.ChangeEvent<{}>, newValue: string) => {
      history.push(newValue);
    };

    return (
      <Paper className={classes.root}>
        <Tabs
          value={history.location.pathname}
          onChange={handleChange}
          indicatorColor="secondary"
          textColor="primary"
          centered
        >
          <Tab value={"/"} label="Home" />
          <Tab value={"/company"} label="Company" />
          <Tab value={"/owner"} label="Owner" />
        </Tabs>
      </Paper>
    );
};

export default withRouter(Navigation);
