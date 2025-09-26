;"use strict";
class Taber {
  constructor() {
    this.tb = document.getElementById('map-nav');

    this.taids = [];
    document.querySelectorAll('.map-nav-bu').forEach(el => {
      let elid = el.getAttribute('id');
      this.taids.push(elid);
      el.addEventListener('click', this.tb_click.bind(this));
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
}
