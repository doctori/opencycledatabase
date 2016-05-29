'use strict';

angular.module('componentService', ['ngResource']).factory('Component', function($resource){
	return $resource('http://192.168.1.15:8080/components/:id',{id:'@ID'},{
		getRange: {
			url: 'http://192.168.1.15:8080/components?page=:page&per_page=:per_page&search_string=:search_string',
			method: 'GET',
			isArray:true,
			params: {
				page:'@page',
				per_page:'@per_page',
				search_string:'@search_string'
			}
		}
	});
}); 