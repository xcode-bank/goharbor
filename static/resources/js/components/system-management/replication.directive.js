(function() {
  
  'use strict';
  
  angular
    .module('harbor.system.management')
    .directive('replication', replication);
  
  ReplicationController.$inject = ['$scope', 'ListReplicationPolicyService', 'ToggleReplicationPolicyService', '$filter', 'trFilter'];
  
  function ReplicationController($scope, ListReplicationPolicyService, ToggleReplicationPolicyService, $filter, trFilter) {
    
    $scope.subsSubPane = 276;
    
    var vm = this;
    vm.retrieve = retrieve;
    vm.search = search;
    vm.togglePolicy = togglePolicy;
    vm.editReplication = editReplication;
    vm.retrieve();
    
    function search() {
      vm.retrieve();
    }
    
    function retrieve() {
      ListReplicationPolicyService('', '', vm.replicationName)
        .success(listReplicationPolicySuccess)
        .error(listReplicationPolicyFailed);
    }
    
    function listReplicationPolicySuccess(data, status) {
      vm.replications = data || [];
    }
    
    function listReplicationPolicyFailed(data, status) {
      $scope.$emit('modalTitle', $filter('tr')('error'));
      $scope.$emit('modalMessage', $filter('tr')('failed_to_list_replication'));
      $scope.$emit('raiseError', true);
      console.log('Failed to list replication policy.');
    }
    
    function togglePolicy(policyId, enabled) {
      ToggleReplicationPolicyService(policyId, enabled)
        .success(toggleReplicationPolicySuccess)
        .error(toggleReplicationPolicyFailed);
    }
    
    function toggleReplicationPolicySuccess(data, status) {
      console.log('Successful toggle replication policy.');
      vm.retrieve();
    }
    
    function toggleReplicationPolicyFailed(data, status) {
      $scope.$emit('modalTitle', $filter('tr')('error'));
      $scope.$emit('modalMessage', $filter('tr')('failed_to_toggle_policy'));
      $scope.$emit('raiseError', true);
      console.log('Failed to toggle replication policy.');
    }
    
    function editReplication(policyId) {
      vm.action = 'EDIT';
      vm.policyId = policyId;
    }
  }
  
  function replication() {
    var directive = {
      'restrict': 'E',
      'templateUrl': '/static/resources/js/components/system-management/replication.directive.html',
      'scope': true,
      'controller': ReplicationController,
      'controllerAs': 'vm',
      'bindToController': true
    };
    return directive;
  }
  
})();