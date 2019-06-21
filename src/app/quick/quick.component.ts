import {Component, OnInit} from '@angular/core';

@Component({
  selector: 'app-quick',
  templateUrl: './quick.component.html',
  styleUrls: ['./quick.component.scss']
})
export class QuickComponent implements OnInit {
  step = 0;
  index = 'First-content';

  constructor() {
  }

  ngOnInit() {
  }

  pre(): void {
    this.step -= 1;
    this.changeContent();
  }

  next(): void {
    this.step += 1;
    this.changeContent();
  }

  done(): void {
    console.log('done');
  }

  changeContent(): void {
    switch (this.step) {
      case 0: {
        this.index = 'First-content';
        break;
      }
      case 1: {
        this.index = 'Second-content';
        break;
      }
      case 2: {
        this.index = 'third-content';
        break;
      }
      default: {
        this.index = 'error';
      }
    }

  }
}
