import { IOwner } from "../owner/ownerInterfaces";

export interface ICompany {
	ID: number;
	Name: string;
	Address: string;
	City: string;
	Country: string;
	Email: string;
	PhoneNumber: string;
	Owners: IOwner;
}
