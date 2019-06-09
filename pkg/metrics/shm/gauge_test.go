/*
 * Licensed to the Apache Software Foundation (ASF) under one or more
 * contributor license agreements.  See the NOTICE file distributed with
 * this work for additional information regarding copyright ownership.
 * The ASF licenses this file to You under the Apache License, Version 2.0
 * (the "License"); you may not use this file except in compliance with
 * the License.  You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package shm

import (
	"os"
	"testing"
	"unsafe"

	"sofastack.io/sofa-mosn/pkg/types"
)

func TestGauge(t *testing.T) {
	types.MosnConfigPath = "/tmp"
	zone := InitMetricsZone("TestGauge", 10*1024)
	defer func() {
		zone.Detach()
		Reset()
	}()

	entry, err := defaultZone.alloc("TestGauge")
	if err != nil {
		t.Error(err)
	}

	gauge := ShmGauge(unsafe.Pointer(&entry.value))

	// update
	gauge.Update(5)

	// value
	if gauge.Value() != 5 {
		t.Error("gauge ops failed")
	}

	gauge.Update(123)
	if gauge.Value() != 123 {
		t.Error("gauge ops failed")
	}
	types.MosnConfigPath = types.MosnBasePath + string(os.PathSeparator) + "conf"
}
