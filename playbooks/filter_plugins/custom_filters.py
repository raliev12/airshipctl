# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

import ipaddress
import types


def is_ipv6(ip):
    return True if ipaddress.ip_address(ip).version == 6 else False


def wrap(ip):
    return '[' + ip + ']' if is_ipv6(ip) else ip


def ipv6_wrap(v):
    if isinstance(v, (list, tuple, types.GeneratorType)):
        return [wrap(ip) for ip in v]
    else:
        return wrap(v)


class FilterModule(object):

    def filters(self):
        return {'ipv6wrap': ipv6_wrap}
