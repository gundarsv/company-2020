import * as React from "react"
import {
	Avatar,
	Dialog,
	DialogContent,
	DialogTitle,
	FormControl,
	IconButton,
	Input,
	InputLabel,
	List,
	ListItem,
	ListItemAvatar,
	ListItemSecondaryAction,
	ListItemText,
	MenuItem,
	Select
} from "@material-ui/core";
import {ICompany} from "../company/companyInterfaces";
import DeleteIcon from '@material-ui/icons/Delete';
import PermIdentityIcon from '@material-ui/icons/PermIdentity';
import {IOwner} from "../owner/ownerInterfaces";
import axios from "axios";

interface IOwnerDialogProps {
	company?: ICompany,
	open: boolean,
	handleClose: () => void;
}

const OwnerDialog: React.FC<IOwnerDialogProps> = (props) => {
	const [owners, setOwners] = React.useState<IOwner[]>();
	const [loading, setLoading] = React.useState(true);

	const loadOwners = async () => {
		setLoading(true);

		await axios.get("/api/owner").then(res => {
			if (props.company.Owners === undefined || props.company.Owners === null) {
				const allOwners: IOwner[] = res.data;
				setOwners(allOwners);
			} else {
				const allOwners: IOwner[] = res.data;
				const availableOwners = allOwners.filter((owner: IOwner) => props.company.Owners.some((currentOwner: IOwner) => owner.ID !== currentOwner.ID));
				setOwners(availableOwners);
			}

		});

		setLoading(false);
	};

	React.useEffect(() => {
		loadOwners();
	}, []);

	return (
		<Dialog open={props.open} onClose={props.handleClose}>
			<DialogTitle id="simple-dialog-title">{props.company.Name + " owners"}</DialogTitle>
			<DialogContent>
				<List dense={false}>
					{
						props.company.Owners ? props.company.Owners.map((owner: IOwner) => {
							return <ListItem key={owner.ID}>
								<ListItemAvatar>
									<Avatar>
										<PermIdentityIcon/>
									</Avatar>
								</ListItemAvatar>
								<ListItemText
									primary={owner.FirstName + " " + owner.LastName}
									secondary={owner.Address}
								/>
								<ListItemSecondaryAction>
									<IconButton edge="end" aria-label="delete" onClick={() => {
										console.log(owner.ID)
									}
									}>
										<DeleteIcon/>
									</IconButton>
								</ListItemSecondaryAction>
							</ListItem>
						}) : null
					}
					<ListItem>
						<FormControl fullWidth={true}>
							<InputLabel htmlFor="grouped-select">Add Owner</InputLabel>
							<Select onChange={object => {
								console.log(object.target.value);
							}} input={<Input id="grouped-select"/>}>
								{
									owners ? owners.map((owner: IOwner) => {
										return <MenuItem
											value={owner.ID}>{owner.FirstName + " " + owner.LastName.charAt(0) + "."}</MenuItem>
									}) : null
								}
							</Select>
						</FormControl>
					</ListItem>
				</List>
			</DialogContent>
		</Dialog>
	)
};

export default OwnerDialog;