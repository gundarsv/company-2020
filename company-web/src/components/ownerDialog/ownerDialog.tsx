import * as React from "react"
import {
	Avatar, CircularProgress,
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
import {useSnackbar} from "notistack";
import AddIcon from '@material-ui/icons/Add';

interface IOwnerDialogProps {
	company?: ICompany,
	open: boolean,
	handleClose: () => void;
}

enum AddDelete {
	AddOwner = 0,
	DeleteOwner = 1
}

const OwnerDialog: React.FC<IOwnerDialogProps> = (props) => {
	const [dropdownOwners, setDropdownOwners] = React.useState<IOwner[]>();
	const [allOwners, setAllOwners] = React.useState<IOwner[]>();
	const [companyOwners, setCompanyOwners] = React.useState<IOwner[]>(props.company.Owners);
	const [loading, setLoading] = React.useState(true);
	const { enqueueSnackbar } = useSnackbar();
	const [ownerSelected, setOwnerSelected] = React.useState<number>();

	const loadOwners = async () => {
		setLoading(true);

		await axios.get("/api/owner").then(res => {
			if (props.company.Owners === undefined || props.company.Owners === null) {
				const allOwners: IOwner[] = res.data;
				setAllOwners(allOwners);
				setDropdownOwners(allOwners);
			} else {
				setAllOwners(res.data);
				setDropdownOwners(filterOwners(res.data));
			}
		});
		setLoading(false);
	};

	const filterOwners = (allOwners:IOwner[]): IOwner[] => {
		return allOwners.filter(owner => !props.company.Owners.some(o => owner.ID == o.ID));
	};

	const updateCompanyOwners = (owner: IOwner, addDelete: AddDelete) => {
		let changedDropdown = dropdownOwners;
		let changedOwners = companyOwners;
		switch (addDelete) {
			case AddDelete.AddOwner:
				changedDropdown.splice(changedDropdown.findIndex((o) =>{
					return o.ID === owner.ID;
				}), 1);
				setDropdownOwners(changedDropdown);
				changedOwners = [...changedOwners, owner];
				setCompanyOwners(changedOwners);
				props.company.Owners = changedOwners;
				break;
			case AddDelete.DeleteOwner:
				changedDropdown = [...changedDropdown, owner];
				setDropdownOwners(changedDropdown);
				changedOwners.splice(changedOwners.findIndex((o) =>{
					return o.ID === owner.ID;
				}), 1);
				setCompanyOwners(changedOwners);
				props.company.Owners = changedOwners;
				break;
			default:
				break;
		}
	};

	const deleteOwnerFromCompany = async (owner: IOwner) => {
		await axios.delete("/api/company/" + props.company.ID + "/owner/" + owner.ID).then(response => {
			if (response.status === 200) {
				setLoading(true);
				updateCompanyOwners(owner, AddDelete.DeleteOwner);
				enqueueSnackbar("Owner was removed from company", {variant: "success"});
				setLoading(false);
				return;
			}
			enqueueSnackbar("Owner was not removed from company", {variant: "error"});
			return;
		})
	};
	
	const addOwnerToCompany = async () => {
		if (ownerSelected === undefined)
		{
			enqueueSnackbar("Owner is not selected", {variant: "warning"});
			return;
		}
		await axios.put("/api/company/" + props.company.ID + "/owner/" + ownerSelected).then(response => {
			if (response.status === 200) {
				setLoading(true);
				updateCompanyOwners(getOwnerByID(allOwners, ownerSelected), AddDelete.AddOwner);
				enqueueSnackbar("Owner was added to company", {variant: "success"});
				setLoading(false);
				return;
			}
			enqueueSnackbar("Owner was not added to company", {variant: "error"});
			return;
		})
	};

	const getOwnerByID =(array: IOwner[], ownerID: number):IOwner=> {
		return array.find(o => ownerID == o.ID);
	};

	React.useEffect(() => {
		loadOwners();
	}, []);

	return (
		<Dialog open={props.open} onClose={props.handleClose}>
			<DialogTitle id="simple-dialog-title">{props.company.Name + " owners"}</DialogTitle>
			<DialogContent>
				{
					!loading ? <List dense={false}>
					{
						companyOwners ? companyOwners.map((owner: IOwner) => {
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
										deleteOwnerFromCompany(owner);
									}
									}>
										<DeleteIcon/>
									</IconButton>
								</ListItemSecondaryAction>
							</ListItem>
						}) : null
					}
					<ListItem>
						<FormControl style={{display:"inline"}} fullWidth={true}>
							<InputLabel htmlFor="select">Add owner</InputLabel>
							<Select style={{width: "87%"}} defaultValue="" input={<Input id="select" disabled={dropdownOwners.length == 0} onChange={(object) => setOwnerSelected(parseInt(object.target.value))} />}>
								{
									dropdownOwners ? dropdownOwners.map((owner: IOwner) => {
											const ownerString:string = owner.FirstName + " " + owner.LastName.charAt(0) + ".";
											return <MenuItem key={owner.ID}
												value={owner.ID}>{ownerString}</MenuItem>
										}) : null
								}
							</Select>
							<IconButton disabled={dropdownOwners.length == 0} edge="end" aria-label="delete" onClick={() => addOwnerToCompany()}>
								<AddIcon/>
							</IconButton>
						</FormControl>
					</ListItem>
					</List> : <div style={{ marginTop: 100, marginBottom: 100, textAlign: "center" }}>
						<CircularProgress size={40}/>
					</div>
				}
			</DialogContent>
		</Dialog>
	)
};

export default OwnerDialog;