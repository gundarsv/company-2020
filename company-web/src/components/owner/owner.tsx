import * as React from "react";
import {IOwner} from "./ownerInterfaces";
import axios from "axios";
import MaterialTable, {Column} from "material-table";
import { CircularProgress } from "@material-ui/core";
import { useSnackbar } from "notistack";

interface OwnerTableState {
	columns: Array<Column<IOwner>>;
	data: IOwner[];
}

const Owner: React.FC = () => {
	const { enqueueSnackbar } = useSnackbar();
	const [loading, setLoading] = React.useState(true);
	const [state, setState] = React.useState<OwnerTableState>({
		columns: [
			{ title: "First Name", field: "FirstName" },
			{ title: "Last Name", field: "LastName" },
			{ title: "Address", field: "Address" },
		],
		data: []
	});

	const loadOwners = async () => {
		setLoading(true);

		await axios.get("/api/owner").then(res => {
			state.data = res.data as IOwner[];
			setState(state);
		});

		setLoading(false);
	};

	const addOwner = (newOwner: IOwner): Promise<IOwner> => {
		return new Promise<IOwner>(async (resolve, reject) => {
			if (!isOwnerValid(newOwner)) {
				enqueueSnackbar("Owner not valid", {variant:"warning"});
				return reject();
			}

			await axios.post("/api/owner", newOwner).then(response => {
				if (response.status === 200)
				{
					setState(prevState => {
						const data = [...prevState.data];
						data.push(response.data);
						return { ...prevState, data };
					});
					enqueueSnackbar("Owner added", {variant:"success"});
					return resolve();
				}
				enqueueSnackbar("Owner was not added", {variant:"error"});
				return reject();
			})
		})
	};

	const deleteOwner = (deletedOwner: IOwner): Promise<IOwner> => {
		return new Promise<IOwner>(async (resolve, reject) => {

			await axios.delete("/api/owner/"+ deletedOwner.ID).then(response => {
				if (response.status === 200)
				{
					setState(prevState => {
						const data = [...prevState.data];
						data.splice(data.indexOf(deletedOwner), 1);
						return { ...prevState, data };
					});
					enqueueSnackbar("Owner was deleted", {variant:"success"});
					return resolve();
				}
				enqueueSnackbar("Owner was not deleted", {variant:"error"});
				return reject();
			})
		})
	};

	const updateOwner = (updatedOwner: IOwner, oldOwner: IOwner): Promise<IOwner> => {
		return new Promise<IOwner>(async (resolve, reject) => {

			if (!isOwnerValid(updatedOwner)) {
				enqueueSnackbar("Updated owner is not valid", {variant:"warning"});
				return reject();
			}

			await axios.put("/api/owner/"+ oldOwner.ID, updatedOwner).then(response => {
				if (response.status === 200)
				{
					setState(prevState => {
						const data = [...prevState.data];
						data[data.indexOf(oldOwner)] = response.data;
						return { ...prevState, data };
					});
					enqueueSnackbar("Owner updated", {variant:"success"});
					return resolve();
				}
				enqueueSnackbar("Owner was not updated", {variant:"error"});
				return reject();
			})
		})
	};

	const isOwnerValid = (owner: IOwner): boolean => {
		if (
			owner.Address == null ||
			owner.Address.length <= 0
		) {
			return false;
		}
		else if (
			owner.FirstName == null ||
			owner.FirstName.length <= 0
		) {
			return false;
		}
		else if (
			owner.LastName == null ||
			owner.LastName.length <= 0
		) {
			return false;
		}

		return true;
	};

	React.useEffect(() => {
		loadOwners();
	}, []);

	if (loading) {
		return (
			<div style={{ marginTop: 300, textAlign: "center" }}>
				<CircularProgress size={100}/>
			</div>
		);
	}
	return (
		<MaterialTable
			style={{ marginTop: 40 }}
			title="Owners"
			columns={state.columns}
			data={state.data}
			editable={{
				onRowAdd: addOwner,
				onRowUpdate: updateOwner,
				onRowDelete: deleteOwner,
			}}
		/>
	);
};

export default Owner;
