/**
 * Sample React Native App
 * https://github.com/facebook/react-native
 *
 * @format
 * @flow strict-local
 */

import React, {useEffect, useState} from 'react';
import {
  SafeAreaView,
  StyleSheet,
  ScrollView,
  View,
  Text,
  StatusBar,
  DeviceEventEmitter,
  NativeModules,
  Button,
} from 'react-native';

import {
  Header,
  LearnMoreLinks,
  Colors,
  DebugInstructions,
  ReloadInstructions,
} from 'react-native/Libraries/NewAppScreen';

const Core = NativeModules.CoreM;

const App: () => React$Node = () => {
  const [pay, setPay] = useState({
    Amount: '',
    MerchantMsg: '',
    Id: '',
  });
  useEffect(() => {
    DeviceEventEmitter.addListener('customEventName', function (e) {
      var data = JSON.parse(e);
      setPay({
        Amount: data.Amount,
        MerchantMsg: data.MerchantMsg,
        Id: data.Id,
      });
      // handle event and you will get a value in event object, you can log it here
    });
  }, []);

  const _pay = () => {
    var obj = {
      Id: pay.Id,
      AccountNo: '555',
      PAN: '12545',
      PcID: '4758',
      ReasonCode: '101',
      ReqID: '1111',
      ReturnCode: '101',
      SerialTransaction: '452111',
      TerminalNo: '4444',
      TraceNumber: '777',
      TransactionDate: '2020-02-01',
      TransactionTime: '20:00',
      Amount: pay.Amount,
    };
    Core.Response(JSON.stringify(obj));
    setPay({
      Amount: '',
      Id: '',
      MerchantMsg: '',
    });
  };

  return (
    <>
      <StatusBar barStyle="dark-content" />
      <SafeAreaView>
        <ScrollView
          contentInsetAdjustmentBehavior="automatic"
          style={styles.scrollView}>
          <View style={styles.body}>
            <View style={styles.sectionContainer}>
              <Text style={styles.sectionDescription}>
                Mobile IP: <Text style={styles.highlight}>{Core.GetIp()}</Text>
              </Text>
            </View>
            <View style={styles.sectionContainer}>
              <Text style={styles.sectionDescription}>
                Amount: <Text style={styles.highlight}>{pay.Amount}</Text>
              </Text>
            </View>
            <View style={styles.sectionContainer}>
              <Text style={styles.sectionDescription}>
                MerchantMsg:{' '}
                <Text style={styles.highlight}>{pay.MerchantMsg}</Text>
              </Text>
            </View>
            <View style={styles.sectionContainer}>
              <Button title="پرداخت" onPress={_pay} />
            </View>
          </View>
        </ScrollView>
      </SafeAreaView>
    </>
  );
};

const styles = StyleSheet.create({
  scrollView: {
    backgroundColor: Colors.lighter,
  },
  engine: {
    position: 'absolute',
    right: 0,
  },
  body: {
    backgroundColor: Colors.white,
  },
  sectionContainer: {
    marginTop: 32,
    paddingHorizontal: 24,
  },
  sectionTitle: {
    fontSize: 24,
    fontWeight: '600',
    color: Colors.black,
  },
  sectionDescription: {
    marginTop: 8,
    fontSize: 18,
    fontWeight: '400',
    color: Colors.dark,
  },
  highlight: {
    fontWeight: '700',
  },
  footer: {
    color: Colors.dark,
    fontSize: 12,
    fontWeight: '600',
    padding: 4,
    paddingRight: 12,
    textAlign: 'right',
  },
});

export default App;
