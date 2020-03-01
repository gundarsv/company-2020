import * as React from "react";
import {ICompany} from "./companyInterfaces";
import axios from "axios";
import MaterialTable, {Column} from "material-table";
import {CircularProgress, Dialog, IconButton, Snackbar} from "@material-ui/core";
import { useSnackbar } from "notistack";
import SupervisorAccountIcon from '@material-ui/icons/SupervisorAccount';
import OwnerDialog from "../ownerDialog/ownerDialog";
import Owner from "../owner/owner";

interface CompanyTableState {
	columns: Array<Column<ICompany>>;
	data: ICompany[];
}

interface OwnerDialogState {
	company: ICompany,
	open: boolean,
}

const Company: React.FC = () => {
	const { enqueueSnackbar } = useSnackbar();
	const [ownerDialogState, setOwnerDialogState] = React.useState({
		company: null,
		open: false
	});
	const [loading, setLoading] = React.useState(true);
	const [state, setState] = React.useState<CompanyTableState>({
		columns: [
			{ title: "Name", field: "Name" },
			{ title: "Address", field: "Address" },
			{ title: "City", field: "City" },
			{ title: "Country", field: "Country" },
			{ title: "Email", field: "Email" },
			{ title: "Phone Number", field: "PhoneNumber", }
		],
		data: []
	});

	const loadCompanies = async () => {
		setLoading(true);

		await axios.get("/api/company").then(res => {
			state.data = res.data as ICompany[];
			setState(state);
		});

		setLoading(false);
	};

	const addCompany = (newCompany: ICompany): Promise<ICompany> => {
		return new Promise<ICompany>(async (resolve, reject) => {
			if (!isCompanyValid(newCompany)) {
				enqueueSnackbar("Company not valid", {variant:"warning"});
				return reject();
			}

			await axios.post("/api/company", newCompany).then(response => {
				if (response.status === 200)
				{
					setState(prevState => {
						const data = [...prevState.data];
						data.push(response.data);
						return { ...prevState, data };
					});
					enqueueSnackbar("Company added", {variant:"success"});
					return resolve();
				}
				enqueueSnackbar("Company was not added", {variant:"error"});
				return reject();
			})
		})
	};

	const deleteCompany = (deletedCompany: ICompany): Promise<ICompany> => {
		return new Promise<ICompany>(async (resolve, reject) => {

			await axios.delete("/api/company/"+ deletedCompany.ID).then(response => {
				if (response.status === 200)
				{
					setState(prevState => {
						const data = [...prevState.data];
						data.splice(data.indexOf(deletedCompany), 1);
						return { ...prevState, data };
					});
					enqueueSnackbar("Company was deleted", {variant:"success"});
					return resolve();
				}
				enqueueSnackbar("Company was not deleted", {variant:"error"});
				return reject();
			})
		})
	};

	const updateCompany = (updatedCompany: ICompany, oldCompany: ICompany): Promise<ICompany> => {
		return new Promise<ICompany>(async (resolve, reject) => {

			if (!isCompanyValid(updatedCompany)) {
				enqueueSnackbar("Updated company is not valid", {variant:"warning"});
				return reject();
			}

			await axios.put("/api/company/"+ oldCompany.ID, updatedCompany).then(response => {
				if (response.status === 200)
				{
					setState(prevState => {
						const data = [...prevState.data];
						data[data.indexOf(oldCompany)] = response.data;
						return { ...prevState, data };
					});
					enqueueSnackbar("Company updated", {variant:"success"});
					return resolve();
				}
				enqueueSnackbar("Company was not updated", {variant:"error"});
				return reject();
			})
		})
	};

	const isCompanyValid = (company: ICompany): boolean => {
		if (
			company.Address == null ||
			company.Address.length <= 0
		) {
			return false;
		}
		else if (
			company.Name == null ||
			company.Name.length <= 0
		) {
			return false;
		}
		else if (
			company.City == null ||
			company.City.length <= 0
		) {
			return false;
		}

		else if (
			company.Country == null ||
			company.Country.length <= 0
		) {
			return false;
		}

		return true;
	};

	const handleOwnerDialogClose = () => {
		setOwnerDialogState({company: null, open: false});
	};

	React.useEffect(() => {
		loadCompanies();
	}, []);

	if (loading) {
		return (
			<div style={{ marginTop: 300, textAlign: "center" }}>
				<CircularProgress size={100}/>
			</div>
		);
	}
	return (
		<>
		<MaterialTable
			actions={[
				{
					icon: () => <SupervisorAccountIcon/>,
					tooltip: 'Show Owners',
					onClick: (event, rowData) => {
						setOwnerDialogState({company: rowData as ICompany, open: true})
					}
				}
			]}
			style={{ marginTop: 40 }}
			title="Companies"
			columns={state.columns}
			data={state.data}
			editable={{
				onRowAdd: addCompany,
				onRowUpdate: updateCompany,
				onRowDelete: deleteCompany,
			}}
		/>
			{
				ownerDialogState.open ? <OwnerDialog company={ownerDialogState.company} open={ownerDialogState.open} handleClose={handleOwnerDialogClose}/> : null
			}
		</>
	);
};

export default Company;
