class GridHandler {
  constructor() {
    this.fun = new Funcs();
    this.fun.ready(this.handler.bind(this));
  }

  handler() {
    const cs = document.getElementsByName('gorilla.csrf.Token')[0].value;

    document.querySelectorAll('.rwd-table-container').forEach((cont) => {
      const grida = new Agrida({
        'gr_cont': cont,
        'cs': cs,
        'fun': this.fun
      });

      grida.docon();
    });
  }
}

new GridHandler();
