(function() {
 
  'use strict';
  
  angular
    .module('harbor.layout.index')
    .controller('IndexController', IndexController);
    
  IndexController.$inject = ['$filter', 'trFilter'];
    
  function IndexController($filter, trFilter) {
    var vm = this;
    
    var indexDesc = $filter('tr')('index_desc', []);
    var indexDesc1 = $filter('tr')('index_desc_1', []);
    var indexDesc2 = $filter('tr')('index_desc_2', []);
    var indexDesc3 = $filter('tr')('index_desc_3', []);
    var indexDesc4 = $filter('tr')('index_desc_4', []);
    var indexDesc5 = $filter('tr')('index_desc_5', []);
    
    vm.message = '<p class="page-content text-justify">' +
        indexDesc + 
			'</p>' +
      '<ul>' +
			 '<li class="long-line">▪︎ ' + indexDesc1 + '</li>' +
			 '<li class="long-line">▪︎ ' + indexDesc2 + '</li>' +
			 '<li class="long-line">▪︎ ' + indexDesc3 + '</li>' +
			 '<li class="long-line">▪︎ ' + indexDesc4 + '</li>' +
			 '<li class="long-line">▪︎ ' + indexDesc5 + '</li>' +
			'</ul>';
    
  }  
      
})();