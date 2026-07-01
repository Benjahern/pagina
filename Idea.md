
# Idea principal
Una pagina para poder vender productos mios como de mi papa, por lo que los 2 seremos los administradores que piodemos ir agregando objetos a la tienda en general, en esta tienda online se debe poder manejar inventario de los objetos que vendamos, son un poco de todo, ya que yo vendo una cosa, y mi papa vende ptra, pero por eso modelaremos items generales y items 32d para yo poder manejar mi inventario y hacer mi calculadora de costes.
La calculadora de costes se basa en que cada objeto 3d creado tiene un coste de crearlo, tanto de filamento, horas y en base a esas horas es por el precio de la luz que consume la impresora 3d.
En esta pagina la gente puede crearse su propia cuienta para comprar, por lo que necesitaremos un documento de acuerdo de privacidad y confidencialidad.
Ademas de poder ver a los usuarios de nuestra pagina, debemos poder gestionarlo, de manera que podeamos ver los pedidos que hagan o poder crear nosotros mismos los pedidos, porque existe el caso en el que en vez de comprar por la pagina, me hablen por whatsapp o por facebook y quieran comprar, por lo que podriamos asignarlo al sistema y poder tener una vision general de lo que el usuwario quiere y prepararlo y tener todo centralizado.
aAdemas me guistarua que si inicio sesion con mi cuenta la cuales la de administrador, deberia poder tener un boton que me da la opcion de poner la vista de ususario normal
# Arquitectura
Todas las cosas limportantes o validaciones se hacen en el frontend y backend, para que el frontend no mande informacion que se pueda capturar en medio de una solicitud al backend y el backend maneja validaciones parano perder informacion y tener la seguridad comppleta
## Backend
El backend lo manejaremos en go porque go es mas rapido y eficniente para paginas web, ademas que tiene un seo mejor por lo que tengo entendido

## Frontend
En frontend manejaremos todo con Nuxt ya que es el que conozco y se que manejka bien con el seop de la pagina

## Auth
El sistema de auth debe ser por jwt mediante el backend, debe poder diferenciar entre entrada de administrador y entrada de ususario.

## Pagos
El sistema de pagos debe ser mediante webpay o flow o lo que investiguemos, ya que ne realidad no se como funciona pero debo tener que investigarlo para poder implemtenrlo en e sistema


