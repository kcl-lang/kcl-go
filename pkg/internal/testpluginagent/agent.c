#include <stdint.h>

extern char *testPluginAgentInvoke(char *method, char *args_json,
                                   char *kwargs_json);

uint64_t test_plugin_agent_get_proxy_ptr() {
  return (uint64_t)(testPluginAgentInvoke);
}
