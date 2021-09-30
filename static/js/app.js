            // dovrebbe essere un esempio di module con javascript
            function Prompt() {
                let toast = function (c) {
                  //definisco quelli che saranno gli argomenti di default della funzione
                    const{
                        msg = '',
                        icon = 'success',
                        position = 'top-end',
        
                    } = c;
        
                    const Toast = Swal.mixin({
                        toast: true,
                        title: msg,
                        position: position,
                        icon: icon,
                        showConfirmButton: false,
                        timer: 3000,
                        timerProgressBar: true,
                        didOpen: (toast) => {
                            toast.addEventListener('mouseenter', Swal.stopTimer)
                            toast.addEventListener('mouseleave', Swal.resumeTimer)
                        }
                    })
        
                    Toast.fire({})
                }
        
                let success = function (c) {
                    const {
                        msg = "",
                        title = "",
                        footer = "",
                    } = c
        
                    Swal.fire({
                        icon: 'success',
                        title: title,
                        text: msg,
                        footer: footer,
                    })
        
                }
        
                let error = function (c) {
                    const {
                        msg = "",
                        title = "",
                        footer = "",
                    } = c
        
                    Swal.fire({
                        icon: 'error',
                        title: title,
                        text: msg,
                        footer: footer,
                    })
        
                }
        
               async function custom(c){
                  const {
                      icon = "",
                      msg = "",
                      title = "",
                      showConfirmButton = true,
                    } = c;
                    
                  const {value: result} = await Swal.fire({
                    icon: icon,
                    title: title,
                    html: msg,
                    backdrop: false,
                    focusConfirm: false,
                    showCancelButton: true,
                    showConfirmButton: showConfirmButton,
                    // prima che il modal si apra inizializza datapicker
                    willOpen: () => {
                      if (c.willOpen !== undefined){
                        c.willOpen();
                      }
                    },  
                    preConfirm: () => {
                      return [
                        document.getElementById('start').value,
                        document.getElementById('end').value
                      ]
                    },
                    // dopo che il modal si è aperto, rimuove disabled e appare datapicker
                    didOpen: () => {
                      if (c.didOpen !== undefined){
                        c.didOpen();
                      }
                    }
        
                  })
                  // se ho un risultato, se non è uguale al bottone cancel, se non è vuoto
                  if (result) {
                      if(result.dismiss !== Swal.DismissReason.cancel){
                          if(result.value !== ""){
                            //faccio un call back
                            if(c.callback !== undefined) {
                              c.callback(result);
                            }
                          } else {
                            c.callback(false);
                         }
                      } else {
                        c.callback(false);
                      }
                   }
    
               }
        
                return {
                    toast: toast,
                    success: success,
                    error: error,
                    custom: custom,
                }
            }