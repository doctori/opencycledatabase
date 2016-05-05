'use strict';

angular.module('componentService', ['ngResource']).factory('Component', function($resource){
	return $resource('http://192.168.198.133:8080/components/:id',{id:'@ID'},{
		// get:{isArray:true},
		// put:{method:'PUT', isArray:false},
		// post:{method:'POST', isArray:false},
		// update:{method: 'POST', isArray:false},

	});
}); 