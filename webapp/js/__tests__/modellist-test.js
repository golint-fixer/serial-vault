/*
 * Copyright (C) 2016-2017 Canonical Ltd
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License version 3 as
 * published by the Free Software Foundation.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 *
 */
//'use strict'

import React from 'react';
import ReactDOM from 'react-dom';
import TestUtils from 'react-addons-test-utils';


jest.dontMock('../components/ModelList');
jest.dontMock('../components/ModelRow');
jest.dontMock('../components/Navigation');


describe('model list', function() {
 it('displays the models page with no models', function() {
	 var ModelList = require('../components/ModelList');
   var IntlProvider = require('react-intl').IntlProvider;
   var Messages = require('../components/messages').en;

   // Mock the data retrieval from the API
   var getModels = jest.genMockFunction();
   ModelList.WrappedComponent.prototype.__reactAutoBindMap.getModels = getModels;

	 // Render the component
	 var modelsPage = TestUtils.renderIntoDocument(
     <IntlProvider locale="en" messages={Messages}>
			 <ModelList />
    </IntlProvider>
	 );

	 expect(TestUtils.isCompositeComponent(modelsPage)).toBeTruthy();

	 // Check all the expected elements are rendered
	 var section = TestUtils.findRenderedDOMComponentWithTag(modelsPage, 'section');
	 var h2 = TestUtils.findRenderedDOMComponentWithTag(modelsPage, 'h2');
	 var nav = TestUtils.findRenderedDOMComponentWithTag(modelsPage, 'nav');

	 // Check that the navigation tag is set correctly
	 expect(nav.firstChild.children.length).toBe(3);
	 expect(nav.firstChild.children[1].firstChild.className).toBe('active');
	 expect(nav.firstChild.children[1].firstChild.textContent).toBe('Models');

   // Check the getModels was called
   expect(getModels.mock.calls.length).toBe(1);

   // Check the 'no models' message is rendered
   expect(section.children.length).toBe(4);
   expect(section.lastChild.textContent).toBe('No models found.');
 });

 it('displays the models page with some models', function() {
	 var ModelList = require('../components/ModelList');
   var IntlProvider = require('react-intl').IntlProvider;
   var Messages = require('../components/messages').en;

   // Shallow render the component with the translations
   const intlProvider = new IntlProvider({locale: 'en', messages: Messages}, {});
   const {intl} = intlProvider.getChildContext();
   var shallowRenderer = TestUtils.createRenderer();

   // Set up a fixture for the model data
   var models = [
     {id: 1, 'brand-id': 'Brand1', model: 'Name1', revision: 11},
     {id: 2, 'brand-id': 'Brand2', model: 'Name2', revision: 22},
     {id: 3, 'brand-id': 'Brand3', model: 'Name3', revision: 33}
   ];

   // Mock the data retrieval from the API
   var getModels = jest.genMockFunction();
   ModelList.WrappedComponent.prototype.__reactAutoBindMap.getModels = getModels;

	 // Render the component
	 shallowRenderer.render(
			 <ModelList.WrappedComponent intl={intl} models={models} />
	 );
   var modelsPage = shallowRenderer.getRenderOutput();

   expect(modelsPage.props.children.length).toBe(3);
   var section = modelsPage.props.children[1];
   expect(section.props.children.length).toBe(4);

   // Check that the table is rendered correctly
   var table = section.props.children[3].props.children;
   var tbody = table.props.children[1]
   expect(tbody.props.children.length).toBe(3); // data rows
   var row1 = tbody.props.children[0];

   expect(row1.type.WrappedComponent.displayName).toBe('ModelRow')
   expect(row1.props.model).toBe(models[0])

 });

});