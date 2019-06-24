import {Component, OnInit} from '@angular/core';

@Component({
  selector: 'app-project',
  templateUrl: './project.component.html',
  styleUrls: ['./project.component.scss']
})
export class ProjectComponent implements OnInit {
  listOfName = [{text: 'Joe', value: 'Joe', byDefault: true}, {text: 'Jim', value: 'Jim'}];
  listOfAddress = [{text: 'London', value: 'London', byDefault: true}, {text: 'Sidney', value: 'Sidney'}];
  listOfSearchName = ['Joe']; // You need to change it as well!
  sortName: string | null = null;
  sortValue: string | null = null;
  searchAddress = 'London';
  listOfData: Array<{ name: string; age: number; address: string; [key: string]: string | number }> = [
    {
      name: 'John Brown',
      age: 32,
      address: 'New York No. 1 Lake Park'
    },
    {
      name: 'Jim Green',
      age: 42,
      address: 'London No. 1 Lake Park'
    },
    {
      name: 'Joe Black',
      age: 32,
      address: 'Sidney No. 1 Lake Park'
    },
    {
      name: 'Jim Red',
      age: 32,
      address: 'London No. 2 Lake Park'
    }
  ];
  // You need to change it as well!
  listOfDisplayData: Array<{ name: string; age: number; address: string; [key: string]: string | number }> = [];


  constructor() {
  }

  ngOnInit() {
  }

}
