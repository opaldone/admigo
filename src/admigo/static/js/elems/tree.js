;"use strict";
class Atree {
  constructor(props_in) {
    this.props = props_in;

    this.wo = window.TRP[this.props.cid];

    this.el = '#' + this.wo.el;

    this.jtable = document.querySelector(this.el + ' > .rwd-table-container');

    this.clexp = 'exp';
    this.tr_exp = 'tr-exp';
    this.expander = this.el + ' ' + '.' + this.tr_exp;

    this.exp_button = null;
    this.col_button = null;

    this.loadRoot();
  }

  showError(err) {
    console.error(err);
  }

  docon() {
    document.querySelectorAll(this.expander).forEach(el => {
      if (this.props.fun.once(el, 'toggle')) return;
      el.addEventListener('click', this.toggle.bind(this));
    });

    this.props.gr.docon();
  }

  init_head_buttons() {
    this.exp_button = document.querySelector(this.el + ' .tree-head-btn.expand');
    this.exp_button.addEventListener('click', this.expand_all.bind(this));
  }

  prepUrl(url_in) {
    let ret = url_in;

    const upar = new URLSearchParams(window.location.search);
    let upar_s = upar.toString();

    if (upar_s.length > 0) {
      ret = url_in + '?';
    }

    ret = ret + upar_s;

    return ret;
  }

  setParentOpened(par) {
    par.classList.add(this.clexp)
  }

  setParentClosed(par) {
    par.classList.remove(this.clexp);
  }

  changeLines() {
    document.querySelectorAll(this.el + " tr").forEach(elem => {
      const h = elem.querySelector('.vline');

      if (!h) return;

      const prid = elem.getAttribute('data-prid');
      const root = elem.getAttribute('data-root');
      let parsel = this.el + ' tr[data-rid="' + prid + '"]' + '[data-root="' + root + '"]';

      const par = document.querySelector(parsel);
      const posele = this.props.fun.position(elem);
      const pospar = this.props.fun.position(par);
      const nh = posele.top - pospar.top;

      h.style.height = nh + 'px';
    });
  }

  createRows(par, trs) {
    let tem = document.createElement('template');
    tem.innerHTML = trs;
    tem.content.querySelectorAll('tr').forEach(ti => {
      par.after(ti);
    });
  }

  removeRows(par) {
    const rid = par.getAttribute('data-rid');
    const root = par.getAttribute('data-root');
    let chsel = this.el + ' tr[data-prid="' + rid + '"]' + '[data-root="' + root + '"]';

    let childs = document.querySelectorAll(chsel);

    if (childs.length == 0) return;

    childs.forEach(ch => {
      this.removeRows(ch);
      ch.remove();
    });
  }

  opa(par, url_in) {
    const url = this.prepUrl(url_in);
    const obj = this.params(par);

    return axios.post(url, obj, {
      headers: {"X-CSRF-Token": this.props.cs}
    });
  }

  parentParams(par) {
    if (!par) {
      return {};
    }

    const lev = par.getAttribute('data-level');
    const did = par.getAttribute('data-id');
    const prid = par.getAttribute('data-rid');
    const root = par.getAttribute('data-root');

    let ret = {
      'lev': lev,
      'par': did,
    }

    if (prid) {
      ret['prid'] = prid;
    }

    if (root) {
      ret['root'] = root;
    }

    return ret;
  }

  params(par) {
    let ret_pars = this.parentParams(par);
    ret_pars['paths'] = this.wo.paths;
    ret_pars['fid'] = this.wo.fid;
    ret_pars['pk'] = this.wo.pk;
    return ret_pars;
  }

  loadRoot() {
    this.opa(null, this.wo.root)
      .then((re) => {
        this.jtable.innerHTML = re.data.cont;
        this.docon();
        this.init_head_buttons();
      })
      .catch((err) => {
        this.showError(err);
      });
  }

  toggle(ev) {
    ev.preventDefault();
    ev.stopPropagation();

    let th = ev.currentTarget;
    let par = this.props.fun.parent(th, 'tr');
    let isOpen = par.classList.contains(this.clexp);

    if (isOpen) {
      this.removeRows(par);
      this.setParentClosed(par);
      this.changeLines();
      return false;
    }

    this.opa(par, this.wo.node)
      .then((re) => {
        this.createRows(par, re.data.cont);
        this.setParentOpened(par);
        this.docon();
        this.changeLines();
      })
      .catch((err) => {
        this.showError(err);
      });

    return false;
  }

  get_closed_now() {
      let clo = [];

      document.querySelectorAll(this.expander).forEach(el => {
        const par = this.props.fun.parent(el, 'tr');

        if (par.classList.contains(this.clexp)) return;

        clo.push({
          'par': par,
          'prom': this.opa(par, this.wo.node)
        });
      });

    return clo;
  }

  expand_all() {
    let closed = this.get_closed_now();

    let imax = closed.length;

    if (imax == 0) {
      this.docon();
      this.changeLines();
      return;
    }

    closed.forEach((cl) => {
      let par = cl.par;
      cl.prom
        .then((re) => {
          this.createRows(par, re.data.cont);
          this.setParentOpened(par)
          imax = imax - 1;
          if (imax == 0) {
            this.expand_all();
          }
        })
        .catch((err) => {
          this.showError(err);
        });
    });
  }
}

class TreeHandler {
  constructor() {
    this.fun = new Funcs();
    this.fun.ready(this.handler.bind(this));
  }

  handler() {
    const cs = document.getElementsByName("gorilla.csrf.Token")[0].value;

    document.querySelectorAll(".tree-cont").forEach(elem => {
      let tbl_cont = elem.querySelector('.rwd-table-container');
      const grida_in = new Agrida({
        'gr_cont': tbl_cont,
        'cs': cs,
        'fun': this.fun
      });

      let cid = elem.id;
      new Atree({
        'cs': cs,
        'cid': cid,
        'fun': this.fun,
        'gr': grida_in
      });
    });
  }
}

new TreeHandler();
