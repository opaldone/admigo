;"use strict";
class Aforms {
  constructor() {
    this.fun = new Funcs();
    this.fun.ready(this.handler.bind(this));
  }

  handler() {
    document.querySelectorAll('input:not([type = hidden])').forEach(el => {
      el.addEventListener('keyup', this.el_keyup.bind(this));
    });
  }

  el_keyup(ev) {
    let ct = ev.currentTarget;
    let par = this.fun.parent(ct, '.form-it-cont');
    let hid = par.querySelector('input[type=hidden]');
    let inv = par.querySelector('.invalid-feedback');
    let ct_val = ct.value;

    if (hid) {
      hid.value = ct_val;
    }

    if (inv && ct_val.length > 0) {
      inv.innerHTML = '';
    }
  }
}

new Aforms();
