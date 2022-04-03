import {connect} from 'react-redux';

import {isAddCardVisible} from 'selectors';

import TodoIssues from './todo_issues';

const mapStateToProps = (state) => ({
    addVisible: isAddCardVisible(state),
});

export default connect(mapStateToProps, null)(TodoIssues);
