(function() {
  
  'use strict';
  
  angular
    .module('harbor.replication', [
      'harbor.services.replication.policy',
      'harbor.services.replication.job'
    ]);
  
})();