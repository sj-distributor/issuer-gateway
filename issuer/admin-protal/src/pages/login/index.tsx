import { FC } from "react"
import Button from "@mui/material/Button"
import FormControl from "@mui/material/FormControl/FormControl"
import Grid from "@mui/material/Grid"
import Stack from "@mui/material/Stack"
import TextField from "@mui/material/TextField"
import FormHelperText from "@mui/material/FormHelperText/FormHelperText"
import LoginRouterGuard from "./components/LoginRouterGuard"
import { useAction } from "./hooks"
import { helpTextStyles, loginBoxStyles } from "./styles"

const Login: FC = () => {
  const { showHelpText, onLogin, onChangeAccount, onChangePassword } =
    useAction()

  return (
    <LoginRouterGuard>
      <Grid container minHeight="100vh">
        <Grid item xs={8}></Grid>
        <Grid item xs={4}>
          <Grid
            container
            justifyContent="center"
            alignItems="center"
            minHeight="100%"
            p={10}
          >
            <Stack spacing={2} sx={loginBoxStyles}>
              <Grid item component="p" textAlign="center">
                登录
              </Grid>
              <FormControl sx={{ mt: 10 }}>
                <Stack spacing={2} component="form">
                  <TextField
                    label="Account"
                    type="text"
                    onChange={onChangeAccount}
                  />
                  <TextField
                    label="Password"
                    type="password"
                    onChange={onChangePassword}
                  />
                  {showHelpText && (
                    <FormHelperText style={helpTextStyles} error>
                      *请输入帐号密码
                    </FormHelperText>
                  )}
                  <Button variant="contained" onClick={onLogin}>
                    登录
                  </Button>
                </Stack>
              </FormControl>
            </Stack>
          </Grid>
        </Grid>
      </Grid>
    </LoginRouterGuard>
  )
}

export default Login
