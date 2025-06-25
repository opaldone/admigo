;"use strict";
class FormTabs {
  constructor() {
    this.fun = new Funcs();
    this.fun.ready(this.handler.bind(this));

    this.tabs = null;
  }

  handler() {
    this.tabs = document.getElementById('form-tabs');

    if (!this.tabs) return;

    document.querySelectorAll(".form-tab").forEach(el => {
      el.addEventListener("click", this.tabClick.bind(this));
    });

    this.actClick();
  }

  actClick() {
    let act_tab = document.querySelector(".form-tab.act");

    if (!act_tab) return;

    this.fun.trigger(act_tab, 'click');
  }

  clearAct() {
    document.querySelectorAll(".form-tab-content").forEach(el => {
      el.classList.remove('act');
    });

    document.querySelectorAll(".form-tab").forEach(el => {
      el.classList.remove('act');
    });
  }

  tabClick(ev) {
    ev.preventDefault();
    ev.stopPropagation();

    let th = ev.currentTarget;
    let tid = th.getAttribute('data-tid');
    let tab = this.tabs = document.getElementById(tid);

    if (!tab) return;

    this.clearAct()

    th.classList.add('act');
    tab.classList.add('act');

    return false;
  }
}

new FormTabs()
