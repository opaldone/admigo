class Dlg {
  constructor(props_in) {
    this.props = props_in;

    this.sh_cl = 'show-dlg-screen';
    this.b_cl = 'dlg-btn';

    this.scr = null;
    this.bd = null;
    this.header = null;
    this.msg = null;
    this.buttons = null;
    this.clobut = null;
    this.prom_good = null;

    this.set_elems();
  }

  set_elems() {
    this.bd = document.body;
    this.scr = document.getElementById('dlg-screen');
    this.header = document.getElementById('dlg-header');
    this.msg = document.getElementById('dlg-msg');
    this.buttons = document.getElementById('dlg-btns');
    this.clobut = document.getElementById('dlg-close');

    this.window_esc = this.window_esc.bind(this);
    this.clobut_click = this.clobut_click.bind(this);

    this.clobut.addEventListener('click', this.clobut_click);
  }

  clobut_click() {
    this.close_dlg();
  }

  set_header(pars) {
    if (!pars.header) return;
    this.header.textContent = pars.header;
  }

  set_msg(pars) {
    if (!pars.msg) return;
    this.msg.textContent = pars.msg;
  }

  btn_click(ev) {
    let th = ev.currentTarget;
    let ret = parseInt(th.getAttribute('data-ret'));
    this.close_dlg(ret);
  }

  set_buttons_click() {
    this.buttons.querySelectorAll('.' + this.b_cl).forEach(el => {
      if (this.props.fun.once(el, 'btn_click')) return;
      el.addEventListener('click', this.btn_click.bind(this));
    });
  }

  set_buttons(pars) {
    if (!pars.buttons) return;

    this.buttons.innerHTML = '';

    pars.buttons.forEach((btn) => {
      let b_el = document.createElement('button');
      b_el.classList.add(this.b_cl);
      if (btn.ret <= 0) b_el.classList.add('cancel');
      b_el.setAttribute('data-ret', btn.ret);
      b_el.textContent = btn.cap;
      this.buttons.appendChild(b_el);
    });

    this.set_buttons_click();
  }

  window_esc(e) {
    if (e.key === 'Escape') {
      this.close_dlg();
    }
  }

  init_dlg(pars) {
    this.set_header(pars);
    this.set_msg(pars);
    this.set_buttons(pars);
    window.addEventListener('keydown', this.window_esc);
  }

  close_dlg(ret_in) {
    let ret = 0;

    if (ret_in != undefined) ret = ret_in;

    window.removeEventListener('keydown', this.window_esc);
    this.bd.classList.remove(this.sh_cl);

    this.prom_good(ret);
  }

  empty_prom() {
    return new Promise((good) => {
      good(0);
    });
  }

  show(pars_in) {
    if (!this.scr) {
      return this.empty_prom();
    }

    if (this.bd.classList.contains(this.sh_cl)) {
      return this.empty_prom();
    }

    let pars = {};
    if (pars_in) pars = pars_in;

    this.init_dlg(pars);

    this.bd.classList.add(this.sh_cl);

    this.scr.focus();
    this.buttons.querySelector("." + this.b_cl).focus();

    return new Promise((good) => {
      this.prom_good = good;
    });
  }
}
