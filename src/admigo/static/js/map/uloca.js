;"use strict";
class Uloca {
  constructor(par_in) {
    this.oin = par_in;
    this.list = document.getElementById('map-users');
    this.user_cnt = document.getElementById('users-cnt');
  }

  docon() {
    this.list.querySelectorAll('li').forEach(el => {
      if (this.oin.fun.once(el, 'li_click')) return;

      el.addEventListener('click', this.li_click.bind(this));
    });
  }

  ref_coo_cont(cid, pos) {
    let liel = document.getElementById(cid);

    if (!liel) return;

    liel.classList.remove('in-se');

    let sid = 'coo-' + cid;
    let el = document.getElementById(sid);

    if (!el) return;

    let msg = 'null';

    if (pos) {
      msg = `${pos.lat};${pos.lng} acc=${pos.acc}`;
    }

    el.textContent = msg;
  }

  ref_ma(cid) {
    let some = this.oin.uslist[cid];

    if (!some) return false;

    this.ref_coo_cont(cid, some.pos);

    if (!some.pos) return false;

    let sp = [some.pos.lat, some.pos.lng];

    if (some.ma) {
      some.ma.setLatLng(sp);
      some.ma.openPopup();
    } else {
      some.ma = L.marker(sp).addTo(this.oin.map);
      let pop = some.nik ? some.nik : some.cid;
      some.ma.bindPopup('<b>' + pop + '</b>').openPopup();
    }

    if (some.ci) {
      some.ci.setLatLng(sp);
      some.ci.setRadius(some.pos.acc);
    } else {
      some.ci = L.circle(sp, {
        color: '#2b5de5',
        fillColor: '#2b5de5',
        fillOpacity: 0.5,
        stroke: false,
        radius: some.pos.acc
      }).addTo(this.oin.map);
    }

    this.oin.map.setView(some.ma.getLatLng(), 17)
  }

  li_click(e) {
    e.preventDefault();
    e.stopPropagation();

    let el = e.currentTarget

    if (el.classList.contains('in-se')) return;

    let cid = el.getAttribute('id');
    el.classList.add('in-se');

    this.oin.req_loca(cid);

    return false;
  }

  ref_cnt() {
    let sc = '';
    let cc = this.list.children.length;

    if (cc > 0) {
      sc = cc;
    }

    this.user_cnt.textContent = sc;
  }

  get_new_li(cid) {
    let si = '<li id="' + cid + '">' +
      '<div class="map-us-cont">' +
      '<div id="nik-' + cid + '"></div>' +
      '<div class="coo-cont" id="coo-' + cid + '">null</div>' +
      '</div>' +
      '<span class="req-loc"><i class="fa-solid fa-spinner"></i></span>' +
      '</li>';

    let tem = document.createElement('template');
    tem.innerHTML = si;

    return tem.content;
  }

  ref_list() {
    const cids = Object.keys(this.oin.uslist);

    if (cids.length == 0) return;

    cids.forEach((cid, _) => {
      let nel = document.getElementById(cid);
      const some = this.oin.uslist[cid];

      if (!nel) {
        nel = this.get_new_li(cid);
        this.list.appendChild(nel);
      }

      let nik = document.getElementById('nik-' + cid);
      nik.textContent = some.nik.length ? some.nik : cid;
    });

    this.docon();

    this.ref_cnt();
  }

  rem_from_list(cid) {
    let nel = document.getElementById(cid);

    if (!nel) return;

    nel.remove();

    this.ref_cnt();
  }
}
