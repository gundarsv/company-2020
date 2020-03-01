import * as React from "react"
import {Dialog} from "@material-ui/core";
import {ICompany} from "../company/companyInterfaces";

interface IOwnerDialogProps {
	company?: ICompany,
	open: boolean,
	handleClose: () => void;
}

const OwnerDialog: React.FC<IOwnerDialogProps> = (props) => {
	return (
		<Dialog open={props.open} onClose={props.handleClose}>
			<h1>"Hello"</h1>
		</Dialog>
	)

}

export default OwnerDialog;