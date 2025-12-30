;"use strict";
class SeList {
  constructor(props_in) {
    this.props = props_in;

    this.form = null;
    if (this.props.able.getAttribute('form-submit') != null) {
      this.form = document.getElementById(this.props.able.getAttribute('form-submit'));
    }

    this.err_cont = null;
    if (this.props.able.getAttribute('data-errid')) {
      this.err_cont = document.getElementById(this.props.able.getAttribute('data-errid'));
    }

    this.list_to_add = null;
    this.items_added = null;
    this.item_cuq = null;
    this.init_list_to_add();

    this.inp = this.props.able.querySelector('input[type=text]');
    this.inp.data = {'tm': null, 'pressed': false};
    this.se_list = this.props.able.querySelector('.se-list');

    this.btn_show = this.props.able.querySelector('.se-list-btn.sh');
    this.btn_cl = this.props.able.querySelector('.se-list-btn.cl');
    this.ex_btn = (this.btn_show != undefined);

    if (this.ex_btn) {
      this.btn_show.addEventListener('click', this.btn_sh_click.bind(this));
      this.btn_cl.addEventListener('click', this.btn_cl_click.bind(this));
    }

    this.inp.addEventListener('keydown', this.inp_keydown.bind(this));
    this.inp.addEventListener('focusout', this.inp_focusout.bind(this));
    this.inp.addEventListener('input', this.inp_input_simple.bind(this));

    if (!this.ex_btn) {
      this.inp.addEventListener('input', this.inp_input.bind(this));
      this.inp.addEventListener('focusin', this.inp_focusin.bind(this));
      this.se_list.addEventListener('focusout', this.se_list_focusout.bind(this));
    }
  }

  mi_item_remove_click(ev) {
    let btn = ev.currentTarget;
    let par_ml_item = this.props.fun.parent(btn, '.ml-item');
    let li_val = btn.getAttribute('data-vuq');

    delete this.items_added[li_val];
    par_ml_item.remove();
  }

  set_mi_item_remove() {
    if (!this.list_to_add) return;
    this.list_to_add.querySelectorAll('.sp-del').forEach(btn => {
      if (this.props.fun.once(btn, 'rem_click')) return;
      if (!btn.getAttribute('data-vuq')) return;
      btn.addEventListener('click', this.mi_item_remove_click.bind(this));
    });
  }

  init_list_to_add() {
    if (!this.props.able.getAttribute('list-to-add')) return;

    this.items_added = {};
    this.list_to_add = document.getElementById(this.props.able.getAttribute('list-to-add'));
    this.item_cuq = this.list_to_add.getAttribute('data-cuq');

    this.list_to_add.querySelectorAll('.ml-item').forEach(liem => {
      let item = liem.querySelector('input[type=hidden]').value
      this.items_added[item] = true;
    });

    this.set_mi_item_remove();
  }

  clear_tm() {
    if (!this.inp.data.tm) return;
    clearTimeout(this.inp.data.tm);
    this.inp.data.tm = null;
  }

  clear_error() {
    if (!this.err_cont) return;
    this.err_cont.innerText = '';
  }

  show_error(er) {
    if (!this.err_cont) return;
    this.err_cont.innerText = er;
  }

  show_se_list() {
    this.se_list.classList.remove('hid');

    if (this.ex_btn) {
      this.btn_cl.classList.remove('hid');
    }
  }

  hide_se_list() {
    this.clear_tm();
    this.se_list.classList.add('hid');

    if (this.ex_btn) {
      this.btn_cl.classList.add('hid');
    }
  }

  clear_search() {
    this.clear_tm();
    this.clear_error();

    this.se_list.innerHTML = '';
    this.hide_se_list();
  }

  get_lis() {
    const ul = this.se_list.querySelector('.search-list-ul');
    if (!ul) return null;

    const list = ul.querySelectorAll('li');
    const lise = ul.querySelector('li.sele');
    const ixsel = Array.prototype.indexOf.call(list, lise);

    return {
      'list': list,
      'lise': lise,
      'ixsel': ixsel
    };
  }

  add_to_list(li) {
    const li_val = li.getAttribute('data-' + this.item_cuq);

    if (this.items_added[li_val]) return;

    const lid = this.list_to_add.id;

    const cp_item = document.getElementById('templ-' + lid);
    let stri = '<li class="ml-item">' +
      cp_item.innerHTML +
    '</li>';

    let tem = document.createElement('template');
    tem.innerHTML = stri;
    let tem_cont = tem.content;

    let inp_hid = tem_cont.querySelector('input[type=hidden]');
    inp_hid.name = lid + '[]';
    inp_hid.value = li_val;
    let sea_el = tem_cont.querySelector('.sea-tmp');
    sea_el.classList.remove('sea-tmp');

    let btn_rem = tem_cont.querySelector('.sp-del');
    btn_rem.setAttribute('data-vuq', li_val);

    this.items_added[li_val] = true;
    this.list_to_add.append(tem_cont);

    this.set_mi_item_remove();
  }

  set_sea(li) {
    for (const [ke, va] of Object.entries(li.dataset)) {
      const cls = '.sea-' + ke;

      document.querySelectorAll(cls).forEach(ele_in => {
        if (ele_in.classList.contains('sea-tmp')) return;

        if (ele_in.classList.contains('sea-inhtml')) {
          if (ele_in.classList.contains('sea-ifemp') && ele_in.textContent.length > 0) return;
          ele_in.textContent = va;
          return;
        }

        if (ele_in.classList.contains('sea-ifemp') && ele_in.value.length > 0) return;

        ele_in.value = va;
        this.props.fun.trigger(ele_in, 'keyup');
      });
    }

    if (this.list_to_add) {
      this.inp.value = '';
    }
  }

  toggle_sele(li) {
    const ul = this.se_list.querySelector('.search-list-ul');
    if (!ul) return null;

    const lise = ul.querySelector('li.sele');
    if (lise) {
      lise.classList.remove('sele');
    }

    li.classList.add('sele');
  }

  li_s_click(li) {
    if (this.list_to_add) {
      this.add_to_list(li);
    }

    this.set_sea(li);

    if (this.ex_btn) {
      this.toggle_sele(li);
      this.hide_se_list();
    } else {
      this.clear_search();
    }

    if (this.form) {
      this.form.submit();
    }
  }

  docon() {
    let lis = this.get_lis();
    if (!lis) return;

    if (!lis.list) return;
    if (lis.list.length == 0) return;

    lis.list.forEach((lit) => {
      if (this.props.fun.once(lit, 'lit_click')) return;
      lit.addEventListener('click', (ev) => {
        ev.stopPropagation();
        ev.preventDefault();

        let ct = ev.currentTarget;
        this.li_s_click(ct);

        return false;
      });
    });
  }

  handler_input() {
    this.inp.data.pressed = true;

    this.clear_search();

    const str = this.inp.value;

    if (str.length == 0) return;

    let url = this.se_list.getAttribute('data-seu');

    url = url + '?' + 'fi=' + str;

    this.inp.data.tm = setTimeout(() => {
      axios.get(url)
        .then((re) => {
          if (!re.data || !re.data.cont) {
            this.clear_search();
            return;
          }

          this.se_list.innerHTML = re.data.cont;
          this.show_se_list();

          this.docon();
        })
        .catch((err) => {
          console.log(err);
          this.show_error(err.response.data.errors.api);
        });
    }, 500);
  }

  inp_input_simple() {
    this.inp.data.pressed = true;
  }

  inp_input() {
    this.handler_input();
  }

  btn_sh_click(ev) {
    ev.stopPropagation();
    ev.preventDefault();

    if (this.inp.data.pressed) {
      this.handler_input();
      this.inp.focus();
      this.inp.data.pressed = false;
      return true;
    }

    const lis = this.get_lis();
    if (!lis) return true;

    this.show_se_list();
    this.inp.focus();
    return true;
  }

  btn_cl_click(ev) {
    ev.stopPropagation();
    ev.preventDefault();

    this.hide_se_list();
    this.inp.focus();

    return true;
  }

  inp_arrow(lis, ev) {
    let sel = lis.lise;

    if (sel == null) {
      lis.list[0].classList.add('sele');
      return false;
    }

    let st = (ev.key == 'ArrowDown' ? 1 : -1);
    let ix_sel = lis.ixsel;
    let ix_nex = ix_sel + st;
    let len = lis.list.length;

    if (ix_nex >= len) return false;
    if (ix_nex < 0) return false;

    sel.classList.remove('sele');
    lis.list[ix_nex].classList.add('sele');

    lis.list[ix_nex].scrollIntoView({ 'behavior': 'smooth', 'block': 'end', 'inline': 'nearest' });

    return false;
  }

  inp_keydown(ev) {
    if (!['Enter', 'ArrowDown', 'ArrowUp'].includes(ev.key)) return true;

    ev.stopPropagation();
    ev.preventDefault();

    const lis = this.get_lis();

    if (!lis || lis.list.length == 0) {
      return false;
    }

    if (['ArrowDown', 'ArrowUp'].includes(ev.key)) {
      return this.inp_arrow(lis, ev);
    }

    if (lis.lise == null) {
      this.li_s_click(lis.list[0]);
      return false;
    }

    this.li_s_click(lis.lise);
    return false;
  }

  inp_focusin() {
    this.props.fun.trigger(this.inp, 'input');
  }

  inp_focusout(ev) {
    if (this.ex_btn) return;

    if (this.se_list == ev.relatedTarget) {
      return true;
    }

    this.clear_search();
  }

  se_list_focusout() {
    this.clear_search();
  }
}

class SeListHandler {
  constructor() {
    this.fun = new Funcs();
    this.fun.ready(this.handler.bind(this));
  }

  handler() {
    document.querySelectorAll('.search-able').forEach(elem => {
      new SeList({
        'fun': this.fun,
        'able': elem
      })
    });
  }
}

new SeListHandler();
