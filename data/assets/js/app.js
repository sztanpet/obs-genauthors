$(document).ready(function() {
  runIfExists('#page-index', setupIndex );
});

function runIfExists( selector, func ) {

  var elem = $( selector );
  if ( elem.length > 0 )
    func( elem );

}

function setupIndex() {
}
