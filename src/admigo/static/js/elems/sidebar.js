;"use strict";
class Sidebar {
  constructor() {
    this.fun = new Funcs();
    this.fun.ready(this.handler.bind(this));

    this.collapa = null;
    this.sb = null;
  }

  handler() {
    this.sb = document.getElementById('sb-el');

    if (!this.sb) return;

    this.collapa = document.getElementById('sb-left');

    document.querySelectorAll(".has-submenu").forEach(el => {
        el.addEventListener("click", this.clickMenu.bind(this));
    });
    this.collapa.addEventListener('click', this.collapa_click.bind(this));

    this.openToAct();
  }

  openToAct() {
    let item = document.querySelector(".sb-item.active");
    let els = [];

    while (item) {
      els.unshift(item);
      item = item.parentNode;
      if (item.classList.contains('zlev')) break;
      if (item.tagName === 'UL') item.classList.add('exp');
      if (item.classList.contains('has-submenu')) item.classList.add('exp');
    }

    setTimeout(() => {
      document.querySelectorAll(".has-submenu").forEach(el => {
        el.classList.remove('init');
      });
    }, 500);
  }

  clickMenu(ev) {
    let trg = ev.target;

    let par = this.fun.parent(trg, 'a');

    if (par || trg.matches('a')) {
      return true;
    }

    ev.preventDefault();
    ev.stopPropagation();

    let th = ev.currentTarget;

    let exp = 'exp';
    let ul = th.getElementsByTagName('ul')[0];
    if (th.classList.contains(exp)) {
      th.classList.remove(exp);
    } else {
      th.classList.add(exp);
    }
    if (ul.classList.contains(exp)) {
      ul.classList.remove(exp);
    } else {
      ul.classList.add(exp);
    }

    return false;
  }

  collapa_click() {
    if (this.sb.classList.contains('opened')) {
      this.sb.classList.add('closing');
      this.sb.classList.remove('opened');
      setTimeout(() => {
        this.sb.classList.remove('closing');
        this.sb.classList.remove('opened');
      }, 300);
      return;
    }

    this.sb.classList.add('opened');
  }
}

new Sidebar();
