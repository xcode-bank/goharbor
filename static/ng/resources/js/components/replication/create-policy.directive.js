(function() {
  
  'use strict';
  
  angular
    .module('harbor.replication')
    .directive('createPolicy', createPolicy);
  
  CreatePolicyController.$inject = ['$scope', 'ListReplicationPolicyService', 'ListDestinationService', 'PingDestinationService', 'CreateReplicationPolicyService', '$location', 'getParameterByName'];
  
  function CreatePolicyController($scope, ListReplicationPolicyService, ListDestinationService, PingDestinationService, CreateReplicationPolicyService, $location, getParameterByName) {
    var vm = this;
    
    //Since can not set value for textarea by using vm
    //use $scope for instead.
    $scope.replication = {};
    $scope.replication.policy = {};
    $scope.replication.destination = {};
    
    var vm0 = $scope.replication.policy;
    var vm1 = $scope.replication.destination;
        
    vm.selectDestination = selectDestination;
    vm.projectId = getParameterByName('project_id', $location.absUrl());
    
    vm.addNew = addNew;
    vm.edit = edit;
    vm.createPolicy = createPolicy;
    vm.pingDestination = pingDestination;
        
    $scope.$watch('vm.destinations', function(current) {
      if(current) {
        console.log('destination:' + angular.toJson(current));
        vm1.selection = current[0]; 
        vm1.endpoint = vm1.selection.endpoint;
        vm1.username = vm1.selection.username;
        vm1.password = vm1.selection.password;
      }
    });
    
    prepareDestination();
        
    $scope.$watch('vm.action', function(current) {
      if(current) {
        console.log('Current action for replication policy:' + current);
        switch(current) {
        case 'ADD_NEW':
          vm.addNew(); break;
        case 'EDIT':
          vm.edit(vm.policyId); break;
        }    
      }
    });
    
    function selectDestination(item) {
      vm1.selection = item;
      vm1.endpoint = item.endpoint;
      vm1.username = item.username;
      vm1.password = item.password;
    }
    
    function prepareDestination() {
      ListDestinationService('')
        .success(listDestinationSuccess)
        .error(listDestinationFailed);
    }

    function addNew() { 
      vm0.name = '';
      vm0.description = '';
      vm0.enabled = true;
    }
    
    function edit(policyId) {
      console.log('Edit policy ID:' + policyId);
      ListReplicationPolicyService(policyId)
        .success(listReplicationPolicySuccess)
        .error(listReplicationPolicyFailed);
    }
    
    function createPolicy(policy) {
      CreateReplicationPolicyService(policy)
        .success(createReplicationPolicySuccess)
        .error(createReplicationPolicyFailed);
    }
    
    function pingDestination() {
      var targetId = vm1.selection.id;
      console.log('Ping target ID:' + targetId);
      PingDestinationService(targetId)
        .success(pingDestinationSuccess)
        .error(pingDestinationFailed);
    }
    
    function listDestinationSuccess(data, status) {
      vm.destinations = data;
    }
    function listDestinationFailed(data, status) {
      console.log('Failed list destination:' + data);
    }
    function listReplicationPolicySuccess(data, status) {
      var replicationPolicy = data[0];
      vm0.name = replicationPolicy.name;
      vm0.description = replicationPolicy.description;
      vm0.enabled = (replicationPolicy.enabled == 1);
      vm.targetId = replicationPolicy.target_id;
    }
    function listReplicationPolicyFailed(data, status) {
      console.log('Failed list replication policy:' + data);
    }
    function createReplicationPolicySuccess(data, status) {
      console.log('Successful create replication policy.');
      vm.clearUp();
    }
    function createReplicationPolicyFailed(data, status) {
      console.log('Failed create replication policy.');
    }
    function pingDestinationSuccess(data, status) {
      alert('Successful ping target:' + data);
    }
    function pingDestinationFailed(data, status) {
      alert('Failed ping target:' + data);
    }
  }
  
  function createPolicy() {
    var directive = {
      'restrict': 'E',
      'templateUrl': '/static/ng/resources/js/components/replication/create-policy.directive.html',
      'scope': {
        'policyId': '@',
        'modalTitle': '@',
        'reload': '&',
        'action': '='
      },
      'link': link,
      'controller': CreatePolicyController,
      'controllerAs': 'vm',
      'bindToController': true
    };
    return directive;
    
    function link(scope, element, attr, ctrl) {
      ctrl.save = save;
      ctrl.clearUp = clearUp;
      function save(form) {
        console.log(angular.toJson(form));
        var postPayload = {
          'projectId': Number(ctrl.projectId),
          'targetId': form.destination.selection.id,
          'name': form.policy.name,
          'enabled': form.policy.enabled ? 1 : 0,
          'description': form.policy.description,
          'cron_str': '',
          'start_time': ''
        };
        ctrl.createPolicy(postPayload, clearUp);
      }
      
      function clearUp() {
        element.find('#createPolicyModal').modal('hide');
        ctrl.reload();
      }
    }
  }
  
})();