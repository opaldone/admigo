;"use strict";
class Taber {
  constructor(fun_in) {
    this.fun = fun_in;
    this.tb = document.getElementById('map-nav');

    this.taids = [];
    document.querySelectorAll('.map-nav-bu').forEach(el => {
      const elid = el.getAttribute('id');
      this.taids.push(elid);
      el.addEventListener('click', this.tb_click.bind(this));
    });

    document.querySelectorAll('.tb-coh, .tb-content').forEach((el) => {
      el.addEventListener('click', () => {
        this._hide_tabs();
      });
    });
  }

  _add_act(tid) {
    document.getElementById(tid).classList.add('act');
  }

  _rem_act(tid) {
    document.getElementById(tid).classList.remove('act');
  }

  _is_shown() {
    return this.tb.classList.contains('sh');
  }

  _is_cls(tid) {
    return this.tb.classList.contains(tid);
  }

  _show_tabs() {
    this.tb.classList.add('sh');
  }

  _hide_tabs() {
    this.tb.classList.remove('sh');

    this.taids.forEach((tid) => {
      this._rem_act(tid);
    });
  }

  _rem_cls(tid) {
    this.tb.classList.remove(tid);
  }

  _add_cls(tid) {
    this.taids.forEach((sid) => {
      if (tid == sid) return;
      this._rem_cls(sid);
      this._rem_act(sid);
    });

    this.tb.classList.add(tid);
    this._add_act(tid);
  }

  tb_click(e) {
    let btn = e.currentTarget;

    let tid = btn.getAttribute('id');

    if (this._is_cls(tid) && this._is_shown()) {
      this._hide_tabs()
      return;
    }

    this._add_cls(tid);
    this._show_tabs();
  }
}
