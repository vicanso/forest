import axios from "axios";

import LoginRegister from "../login_register";
import { USERS_LOGIN } from "../../urls";
import { message } from "antd";
import PropTypes from "prop-types";

class Login extends LoginRegister {
  constructor(props) {
    super(props);
    this.state.mode = this.loginMode;
  }
  async componentWillMount() {
    try {
      const { data } = await axios.get(USERS_LOGIN);
      this.setState({
        token: data.token
      });
    } catch (err) {
      message.error(err.message);
    }
  }
}

Login.propTypes = {
  setUserInfo: PropTypes.func.isRequired
};

export default Login;
