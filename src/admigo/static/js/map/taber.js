;"use strict";
class Taber {
  constructor(fun_in) {
    this.fun = fun_in;
    this.tb = document.getElementById('map-nav');
    this.tb_users = document.getElementById('tb-users');

    this.taids = [];
    document.querySelectorAll('.map-nav-bu').forEach(el => {
      let elid = el.getAttribute('id');
      this.taids.push(elid);
      el.addEventListener('click', this.tb_click.bind(this));
    });

    document.querySelectorAll('.tb-coh, .tb-content').forEach((el) => {
      el.addEventListener('click', () => {
        this.hide_tabs();
      });
    });
  }

  _rem_cls(tid) {
    this.tb.classList.remove(tid);
    this.tb.classList.remove('sh');
    document.getElementById(tid).classList.remove('act');
  }

  _add_cls(tid) {
    this.tb.classList.add(tid);
    this.tb.classList.add('sh');
    document.getElementById(tid).classList.add('act');
  }

  _clear_cls(tid) {
    this.taids.forEach((sid) => {
      if (sid == tid) return;
      this._rem_cls(sid);
    });
  }

  hide_tabs() {
    this._clear_cls('');
  }

  tb_click(e) {
    let btn = e.currentTarget;

    let tid = btn.getAttribute('id');

    if (this.tb.classList.contains(tid)) {
      this._rem_cls(tid);
      return;
    }

    this._clear_cls(tid);
    this._add_cls(tid);
  }

  show_users() {
    if (this.tb.classList.contains('tb-users')) return;

    this.fun.trigger(this.tb_users, 'click');
  }
}
